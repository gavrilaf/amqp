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

	var res float32
	var err error

	switch arg.Op {
	case "+":
		res = arg.Left + arg.Right
	case "-":
		res = arg.Left - arg.Right
	case "*":
		res = arg.Left * arg.Right
	case "/":
		if arg.Right == 0.0 {
			err = errors.New("divide by 0")
		} else {
			res = arg.Left / arg.Right
		}
	default:
		err = errors.New("unkown function")
	}

	return &Answer{Result: res, Req: arg}, err
}
