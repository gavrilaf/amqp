RPC implementation over RabbitMQ

# Installation

Install the standard protocol buffer implementation from https://github.com/google/protobuf.
Install the gogo golang protobuf implementation from https://github.com/gogo/protobuf.

Install mq-rpc binary:

`go get github.com/gavrilaf/amqp/protoc-gen-mqrpc`

# Run example

Run RabbitMQ locally or from Docker image

`docker run --hostname my-rabbit -p 5672:5672 --name rabbit22 rabbitmq:3`

Run amqp/examples/rpc-test-srv

`go run server.go messages.pb.go messages_mqrpc.gen.go`

`go run client.go messages.pb.go messages_mqrpc.gen.go`

# Using library

Write service definition:
```protobuf
service TestService {
  rpc Ping (Empty) returns (ServerStatus);
  rpc CreateUser(User) returns (ResourceID);
  rpc CreateAccount(Account) returns (ResourceID);
  rpc FindAccount(ResourceID) returns (Account);
}

message Empty {}

message ResourceID {
  string ID = 1;
}

message ServerStatus {
  int32 status = 1; 
}
......
```
Generate types & service:
```
protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogofast_out=. *.proto

protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --mqrpc_out=. *.proto
```

Impement server:
``` golang
func main() {
  srv, err := rpc.CreateServer("amqp://localhost:5672", "rpc-rabbit-worker")
  RunServer(srv, &srvHandler{})
}

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
......
```

Connect to server and use service methods:
``` golang
func main() {
  conn, err := rpc.Connect(rpc.ClientConfig{Url: "amqp://localhost:5672", 
                                            ServerQueue: "rpc-rabbit-worker", 
                                            Timeout: time.Second})
  
  client := NewTestServiceClient(conn)
  
  status, err := client.Ping(&Empty{})
	
  id1, err := client.CreateUser(&User{Username: "username", PasswordHash: "111", Device: nil})
......
```




