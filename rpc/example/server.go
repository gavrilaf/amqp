package main

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/amqp/rpc"
	"github.com/satori/go.uuid"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s\n", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	fmt.Printf("Starting AMQP RPC server\n")

	handler := serverImpl{}

	srv, err := rpc.CreateRpcServer("amqp://localhost:5672", "rpc-rabbit-worker")
	if err != nil {
		fmt.Printf("Error starting server %s\n", err)
		panic(err)
	}

	srv.Serve(func(funcID int32, arg []byte) ([]byte, error) {
		switch Functions(funcID) {
		case Functions_Ping:
			return handler.HandlePing(arg)
		case Functions_CreateUser:
			return handler.HandleCreateUser(arg)
		case Functions_CreateAccount:
			return handler.HandleCreateAccount(arg)
		default:
			return nil, errors.New("unknown function")
		}
	})
}

/////////////////////////////////////////////////////////////////////////////////////

type Server interface {
	Ping() (*ServerPingResponse, error)
	CreateUser(user CreateUserRequest) (*IDResponse, error)
	CreateAccount(account CreateAccountRequest) (*IDResponse, error)
}

type serverImpl struct{}

func (p serverImpl) Ping() (*ServerPingResponse, error) {
	fmt.Printf("Ping call\n")
	return &ServerPingResponse{Status: 2}, nil
}

func (p serverImpl) CreateUser(user CreateUserRequest) (*IDResponse, error) {
	fmt.Printf("Create user call: %v\n", spew.Sdump(user))
	return &IDResponse{ID: uuid.NewV4().String()}, nil
}

func (p serverImpl) CreateAccount(account CreateAccountRequest) (*IDResponse, error) {
	fmt.Printf("Create account call: %v\n", spew.Sdump(account))
	return nil, errors.New("Creating account error")
}

/////////////////////////////////////////////////////////////////////////////////////

func (p serverImpl) HandlePing(arg []byte) ([]byte, error) {
	resp, err := p.Ping()
	if err != nil {
		return nil, err
	}

	return resp.Marshal()
}

func (p serverImpl) HandleCreateUser(arg []byte) ([]byte, error) {
	var req CreateUserRequest
	err := req.Unmarshal(arg)
	if err != nil {
		return nil, err
	}

	resp, err := p.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return resp.Marshal()
}

func (p serverImpl) HandleCreateAccount(arg []byte) ([]byte, error) {
	var req CreateAccountRequest
	err := req.Unmarshal(arg)
	if err != nil {
		return nil, err
	}

	resp, err := p.CreateAccount(req)
	if err != nil {
		return nil, err
	}

	return resp.Marshal()
}
