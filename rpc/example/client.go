package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gavrilaf/amqp/rpc"
)

func main() {
	fmt.Printf("Starting AMQP RPC client\n")

	rpc, err := rpc.Connect("amqp://localhost:5672", "rpc-rabbit-worker")
	defer rpc.Close()

	if err != nil {
		panic(fmt.Errorf("Dial error: %v\n", err))
	}

	client := rpcClient{rpc: rpc}

	fmt.Printf("Call Ping\n")
	resp, err := client.Ping(Empty{})
	if err != nil {
		fmt.Printf("Ping error: %v\n", err)
	} else {
		fmt.Printf("Ping result: %v\n", spew.Sdump(resp))
	}

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
	idResp, err := client.CreateUser(req1)
	if err != nil {
		fmt.Printf("CreateUser error: %v\n", err)
	} else {
		fmt.Printf("CreateUser ok, %s\n", spew.Sdump(idResp))
	}

	req2 := CreateAccountRequest{Currency: "USD"}

	fmt.Printf("Call CreateAccount(%s)\n", spew.Sdump(req2))
	idResp, err = client.CreateAccount(req2)
	if err != nil {
		fmt.Printf("CreateAccount error: %v\n", err)
	} else {
		fmt.Printf("CreateAccount ok, %s\n", spew.Sdump(idResp))
	}

	select {}
}

type rpcClient struct {
	rpc *rpc.RpcClient
}

func (client rpcClient) Ping(p Empty) (*ServerPingResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := client.rpc.Call(rpc.Request{FuncID: int32(Functions_Ping), Body: request})
	if err != nil {
		return nil, err
	}

	var resp ServerPingResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}

func (client rpcClient) CreateUser(p CreateUserRequest) (*IDResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := client.rpc.Call(rpc.Request{FuncID: int32(Functions_CreateUser), Body: request})
	if err != nil {
		return nil, err
	}

	var resp IDResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}

func (client rpcClient) CreateAccount(p CreateAccountRequest) (*IDResponse, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := client.rpc.Call(rpc.Request{FuncID: int32(Functions_CreateAccount), Body: request})
	if err != nil {
		return nil, err
	}

	var resp IDResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}
