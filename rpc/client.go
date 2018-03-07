package rpc

import (
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

var ErrTimeout = errors.New("timeout")

type Client interface {
	Close()
	RemoteCall(p Request) ([]byte, error)
}

type ClientConfig struct {
	Url         string
	ServerQueue string
	Timeout     time.Duration
}

func Connect(cfg ClientConfig) (Client, error) {
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
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

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	client := newClient(cfg.ServerQueue, conn, channel, &queue, cfg.Timeout)
	go client.handleDeliveries(msgs)

	return client, nil
}

///////////////////////////////////////////////////////////////////////////////////

type clientImpl struct {
	serverQueue string
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       *amqp.Queue

	timeout time.Duration
	calls   map[string]*pendingCall

	done chan bool
}

type pendingCall struct {
	done chan bool
	data []byte
}

func newClient(serverQueue string, conn *amqp.Connection, channel *amqp.Channel, queue *amqp.Queue, timeout time.Duration) *clientImpl {
	return &clientImpl{
		serverQueue: serverQueue,
		conn:        conn,
		channel:     channel,
		queue:       queue,
		timeout:     timeout,
		calls:       make(map[string]*pendingCall),
		done:        make(chan bool)}
}

func (client *clientImpl) Close() {
	if client == nil {
		return
	}

	client.done <- true

	if client.channel != nil {
		client.channel.Close()
	}

	if client.conn != nil {
		client.conn.Close()
	}
}

func (client *clientImpl) RemoteCall(p Request) ([]byte, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	expiration := fmt.Sprintf("%d", client.timeout)
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

	call := &pendingCall{done: make(chan bool)}
	client.calls[corrId] = call

	var respData []byte
	var respError error = ErrTimeout

	select {
	case <-call.done:
		var resp Response
		respError = resp.Unmarshal(call.data)
		if respError == nil {
			if resp.IsSuccess {
				respData = resp.Body
			} else {
				respError = errors.New(resp.ErrText)
			}
		}

	case <-time.After(client.timeout):
		break
	}

	delete(client.calls, corrId)

	return respData, respError
}

func (client *clientImpl) handleDeliveries(msgs <-chan amqp.Delivery) {
	finish := false
	for !finish {
		select {
		case msg := <-msgs:
			call, ok := client.calls[msg.CorrelationId]
			if ok {
				call.data = msg.Body
				call.done <- true
			}
		case <-client.done:
			finish = true
		}
	}
}

func newCorrId() string {
	return uuid.NewV4().String()
}
