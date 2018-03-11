package main

import (
	"fmt"
	"github.com/gavrilaf/amqp/rpc"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Printf("Starting AMQP RPC calculator client\n")

	requests := []Request{
		Request{Left: 10, Right: 1, Op: "+"},
		Request{Left: 1000, Right: 13.4, Op: "*"},
		Request{Left: 98, Right: 0, Op: "/"},
		Request{Left: 123, Right: 1000, Op: "-"},
		Request{Left: 10.23, Right: 13.13, Op: "+"},
		Request{Left: 10, Right: 10, Op: "*"},
		Request{Left: 10, Right: 189, Op: "+"},
		Request{Left: 10, Right: 10, Op: "/"},
		Request{Left: 10, Right: 0, Op: "/"},
		Request{Left: 2, Right: 2, Op: "*"},
		Request{Left: 780, Right: 123, Op: "+"},
		Request{Left: 9000, Right: 9000, Op: "-"},
		Request{Left: 9000, Right: 9000, Op: "()"},
	}

	srv, err := rpc.Connect(rpc.ClientConfig{Url: "amqp://localhost:5672", ServerQueue: "rpc-rabbit-worker", Timeout: time.Second})
	failOnError(err, "Connect")

	client := NewCalcClient(srv)

	for _, req := range requests {
		answer, err := client.Eval(&req)
		fmt.Printf("%s : (%s: %v) \n", req.String(), answer.String(), err)
	}
}
