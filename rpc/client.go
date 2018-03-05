package rpc

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

var ErrTimeout = errors.New("timeout")

func Connect(url string, serverQueue string) (*RpcClient, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	client := newRpcClient(serverQueue, conn, ch, &q)
	go client.handleDeliveries(msgs)

	return client, nil
}

///////////////////////////////////////////////////////////////////////////////////

type RpcClient struct {
	serverQueue string
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       *amqp.Queue

	calls map[string]*rpcCall
}

type rpcCall struct {
	done chan bool
	data []byte
}

func newRpcClient(serverQueue string, conn *amqp.Connection, ch *amqp.Channel, q *amqp.Queue) *RpcClient {
	client := &RpcClient{serverQueue: serverQueue, conn: conn, channel: ch, queue: q, calls: make(map[string]*rpcCall)}
	return client
}

func (client *RpcClient) Close() {
	if client != nil && client.channel != nil {
		client.channel.Close()
	}
	if client != nil && client.conn != nil {
		client.conn.Close()
	}
}

func (client *RpcClient) Call(p CallDesc) ([]byte, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	expiration := fmt.Sprintf("%d", 1000)
	corrId := newCorrId()
	err = client.channel.Publish(
		"",                 // exchange
		client.serverQueue, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: corrId,
			ReplyTo:       client.queue.Name,
			Body:          request,
			Expiration:    expiration,
		})
	if err != nil {
		return nil, err
	}

	rpcCall := &rpcCall{done: make(chan bool)}
	client.calls[corrId] = rpcCall

	var resData []byte
	var resError error

	select {
	case <-rpcCall.done:
		resData = rpcCall.data

	case <-time.After(time.Millisecond * time.Duration(3)):
		return nil, ErrTimeout
	}

	delete(client.calls, corrId)

	return resData, resError
}

func (client *RpcClient) handleDeliveries(msgs <-chan amqp.Delivery) {
	for d := range msgs {

		call, ok := client.calls[d.CorrelationId]

		if !ok {
			continue
		}

		call.data = d.Body
		call.done <- true
	}
}

func newCorrId() string {
	return uuid.NewV4().String()
}
