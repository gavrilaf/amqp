package rpc

import (
	"fmt"
	"github.com/streadway/amqp"
)

type CallHandler func(funcID int32, args []byte) ([]byte, error)

type Server interface {
	Close()
	Serve(handler CallHandler)
}

func CreateServer(url string, queueName string) (Server, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	return &serverImpl{conn: conn, channel: ch, msgs: msgs, done: make(chan bool)}, nil
}

type serverImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	done    chan bool
}

func (srv *serverImpl) Close() {
	if srv == nil {
		return
	}

	srv.done <- true

	if srv.channel != nil {
		srv.channel.Close()
	}

	if srv.conn != nil {
		srv.conn.Close()
	}
}

func (srv *serverImpl) Serve(handler CallHandler) {
	finish := false
	for !finish {
		select {
		case msg := <-srv.msgs:
			go srv.callHandler(handler, msg)
		case <-srv.done:
			finish = true
		}
	}
}

func (srv *serverImpl) callHandler(handler CallHandler, msg amqp.Delivery) {
	var req Request
	err := req.Unmarshal(msg.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed unmarshal request: %v", err))
	}

	var resp Response
	data, err := handler(req.FuncID, req.Body)
	if err != nil {
		resp.IsSuccess = false
		resp.ErrText = err.Error()
	} else {
		resp.IsSuccess = true
		resp.Body = data
	}

	respData, err := resp.Marshal()
	if err != nil {
		panic(fmt.Sprintf("Failed marshall responce: %v", err))
	}

	err = srv.channel.Publish(
		"",          // exchange
		msg.ReplyTo, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: msg.CorrelationId,
			Body:          respData,
		})

	if err != nil {
		panic(fmt.Sprintf("Failed to publish a message: %v", err))
	}

	msg.Ack(false)
}
