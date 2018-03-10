package main

import (
	//"fmt"
	//	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	//"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	p := NewMqRpc()
	resp := command.GeneratePlugin(req, p, "_mqrpc.gen.go")
	command.Write(resp)
}
