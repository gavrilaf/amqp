package main

import (
	"errors"
	"fmt"
	"github.com/gavrilaf/amqp/rpc"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Printf("Starting calc server\n")

	srv, err := rpc.CreateServer("amqp://localhost:5672", "rpc-rabbit-worker")
	failOnError(err, "RabbitMQ connection")

	RunServer(srv, &calcHandler{})
}

/////////////////////////////////////////////////////////////////////////////////////

type calcHandler struct{}

func (p calcHandler) Eval(arg *Request) (*Answer, error) {
	fmt.Printf("Eval: %s\n", arg.String())

	switch arg.Op {
	case "+":
		return &Answer{Result: arg.Left + arg.Right}, nil
	case "-":
		return &Answer{Result: arg.Left - arg.Right}, nil
	case "*":
		return &Answer{Result: arg.Left * arg.Right}, nil
	case "/":
		if arg.Right == 0.0 {
			return nil, errors.New("divide by 0")
		} else {
			return &Answer{Result: arg.Left / arg.Right}, nil
		}
	default:
		return nil, errors.New("unkown function")
	}
}
