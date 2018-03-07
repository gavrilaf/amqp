package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/amqp/rpc"
	"time"
)

func main() {
	fmt.Printf("Starting AMQP RPC client\n")

	client, err := rpc.Connect(rpc.ClientConfig{Url: "amqp://localhost:5672", ServerQueue: "rpc-rabbit-worker", Timeout: time.Second * 1})
	defer client.Close()

	if err != nil {
		panic(fmt.Errorf("Dial error: %v\n", err))
	}

	bridge := Bridge{client: client}

	////////////////////////////////////////////////////////////////////
	fmt.Printf("Call Ping\n")
	resp, err := bridge.Ping(Empty{})
	if err != nil {
		fmt.Printf("Ping error: %v\n", err)
	} else {
		fmt.Printf("Ping result: %v\n", spew.Sdump(resp))
	}

	////////////////////////////////////////////////////////////////////
	req1 := CreateUserRequest{
		Username:     "user",
		PasswordHash: "123456",
		Device: &Device{
			ID:     "device-1",
			Name:   "Test device",
			Locale: "ru",
			Lang:   "es"},
	}
	fmt.Printf("Call CreateUser(%s)\n", spew.Sdump(req1))

	idResp, err := bridge.CreateUser(req1)
	if err != nil {
		fmt.Printf("CreateUser error: %v\n", err)
	} else {
		fmt.Printf("CreateUser ok, %s\n", spew.Sdump(idResp))
	}

	////////////////////////////////////////////////////////////////////
	req2 := CreateAccountRequest{Currency: "USD"}
	fmt.Printf("Call CreateAccount(%s)\n", spew.Sdump(req2))

	idResp, err = bridge.CreateAccount(req2)
	if err != nil {
		fmt.Printf("CreateAccount error: %v\n", err)
	} else {
		fmt.Printf("CreateAccount ok, %s\n", spew.Sdump(idResp))
	}
}

type Bridge struct {
	client rpc.Client
}

func (bridge Bridge) Ping(p Empty) (*ServerPingResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := bridge.client.RemoteCall(rpc.Request{FuncID: int32(Functions_Ping), Body: request})
	if err != nil {
		return nil, err
	}

	var resp ServerPingResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}

func (bridge Bridge) CreateUser(p CreateUserRequest) (*IDResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := bridge.client.RemoteCall(rpc.Request{FuncID: int32(Functions_CreateUser), Body: request})
	if err != nil {
		return nil, err
	}

	var resp IDResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}

func (bridge Bridge) CreateAccount(p CreateAccountRequest) (*IDResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := bridge.client.RemoteCall(rpc.Request{FuncID: int32(Functions_CreateAccount), Body: request})
	if err != nil {
		return nil, err
	}

	var resp IDResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}
