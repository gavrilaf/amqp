package main

import (
	"errors"
	"fmt"
	"github.com/gavrilaf/amqp/rpc"
	"github.com/satori/go.uuid"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Printf("Starting AMQP RPC server\n")

	srv, err := rpc.CreateServer("amqp://localhost:5672", "rpc-rabbit-worker")
	failOnError(err, "RabbitMQ connection")

	defer srv.Close()

	RunServer(srv, &srvHandler{})
}

/////////////////////////////////////////////////////////////////////////////////////

type srvHandler struct{}

func (p srvHandler) Ping(arg *Empty) (*ServerStatus, error) {
	fmt.Printf("Ping\n")
	return &ServerStatus{Status: 2}, nil
}

func (p srvHandler) CreateUser(user *User) (*ResourceID, error) {
	fmt.Printf("CreateUser: %v\n", user.String())
	return &ResourceID{ID: uuid.NewV4().String()}, nil
}

func (p srvHandler) CreateAccount(acc *Account) (*ResourceID, error) {
	fmt.Printf("CreateAccount: %v\n", acc.String())
	return nil, errors.New("Creating account error")
}

func (p srvHandler) FindAccount(acc *ResourceID) (*Account, error) {
	fmt.Printf("FindAccount: %v\n", acc.String())

	// Simulate slow call
	time.Sleep(time.Second * 2)
	return nil, errors.New("Find account error")
}
