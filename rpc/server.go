package rpc

import (
	//"fmt"
	"github.com/streadway/amqp"
)

type RpcServer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	msgs    <-chan amqp.Delivery
	done    chan bool
}

func CreateRpcServer(url string, queueName string) (*RpcServer, error) {
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

	return &RpcServer{conn: conn, channel: ch, msgs: msgs, done: make(chan bool)}, nil
}

type RpcHandler func(funcID int32, args []byte) ([]byte, error)

func (srv *RpcServer) Serve(handler RpcHandler) {
	for msg := range srv.msgs {
		var desc CallDesc
		err := desc.Unmarshal(msg.Body)
		if err != nil {
			panic("Something wrong!!!")
		}

		resp, err := handler(desc.FuncID, desc.Msg)
		if err != nil {
			panic("Do something")
		}

		err = srv.channel.Publish(
			"",          // exchange
			msg.ReplyTo, // routing key
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType:   "application/octet-stream",
				CorrelationId: msg.CorrelationId,
				Body:          resp,
			})

		if err != nil {
			panic("error: Failed to publish a message")
		}

		msg.Ack(false)
	}
}
