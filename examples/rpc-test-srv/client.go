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
	fmt.Printf("Starting AMQP RPC client\n")

	srv, err := rpc.Connect(rpc.ClientConfig{Url: "amqp://localhost:5672", ServerQueue: "rpc-rabbit-worker", Timeout: time.Second})
	failOnError(err, "Connect")

	client := NewTestServiceClient(srv)
	defer client.Close()

	status, err := client.Ping(&Empty{})
	fmt.Printf("Ping, result = %s, error = %v\n", status.String(), err)

	id1, err := client.CreateUser(&User{Username: "username", PasswordHash: "111", Device: nil})
	fmt.Printf("CreateUser, result = %s, error = %v\n", id1.String(), err)

	id2, err := client.CreateAccount(&Account{Currency: "USD"})
	fmt.Printf("CreateAccount, result = %s, error = %v\n", id2.String(), err)

	acc, err := client.FindAccount(&ResourceID{ID: "12345"})
	fmt.Printf("FindAccount, result = %s, error = %v\n", acc.String(), err)
}
