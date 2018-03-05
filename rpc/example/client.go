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

	req := CreateUserRequest{
		Username:     "user",
		PasswordHash: "123456",
		Device: &Device{
			ID:     "device-1",
			Name:   "Test device",
			Locale: "ru",
			Lang:   "es"},
	}

	fmt.Printf("Call CreateUser(%s)\n", spew.Sdump(req))
	_, err = client.CreateUser(req)
	if err != nil {
		fmt.Printf("CreateUser error: %v\n", err)
	} else {
		fmt.Printf("CreateUser ok\n")
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

	respData, err := client.rpc.Call(rpc.CallDesc{FuncID: int32(Functions_Ping), Msg: request})
	if err != nil {
		return nil, err
	}

	var resp ServerPingResponse
	err = resp.Unmarshal(respData)

	return &resp, err
}

func (client rpcClient) CreateUser(p CreateUserRequest) (*Empty, error) {
	request, err := p.Marshal()
	if err != nil {
		return nil, err
	}

	respData, err := client.rpc.Call(rpc.CallDesc{FuncID: int32(Functions_CreateUser), Msg: request})
	if err != nil {
		return nil, err
	}

	var resp Empty
	err = resp.Unmarshal(respData)

	return &resp, err
}
