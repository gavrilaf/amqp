generate:
	# Generate examles & test output
	protoc -I=./examples/rpc-calc -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./examples/rpc-calc/ ./examples/rpc-calc/*.proto
	
	protoc -I=./examples/rpc-calc \
		-I=$$GOPATH/src -I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./examples/rpc-calc/ ./examples/rpc-calc/*.proto

	protoc -I=./examples/rpc-test-srv -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./examples/rpc-test-srv/ ./examples/rpc-test-srv/*.proto
	
	protoc -I=./examples/rpc-test-srv -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./examples/rpc-test-srv/ ./examples/rpc-test-srv/*.proto

	protoc -I=./rpc/test -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./rpc/test/ ./rpc/test/*.proto

	protoc -I=./rpc/test -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./rpc/test/ ./rpc/test/*.proto

test:
	# 
	protoc -I=./rpc/test -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--gogofast_out=./rpc/test/ ./rpc/test/*.proto

	protoc -I=./rpc/test -I=$$GOPATH/src \
		-I=$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		--mqrpc_out=./rpc/test/ ./rpc/test/*.proto

	go test ./rpc/test