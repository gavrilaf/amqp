RPC implementation over RabbitMQ

Run example

Run RabbitMQ
`docker run --hostname my-rabbit -p 5672:5672 --name rabbit22 rabbitmq:3`

Run server (amqp/rpc/example)
`go run server.go messages.pb.go`

Run client (amqp/rpc/example)
`go run client.go messages.pb.go`


