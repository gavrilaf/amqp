RPC implementation over RabbitMQ

Run example

Firstly you have to run RabbitMQ locally or from Docker image

`docker run --hostname my-rabbit -p 5672:5672 --name rabbit22 rabbitmq:3`

Run amqp/examples/rpc-test-srv

`go run server.go messages.pb.go messages_mqrpc.gen.go`

`go run client.go messages.pb.go messages_mqrpc.gen.go`


