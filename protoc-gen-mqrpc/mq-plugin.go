package main

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"strings"
)

type mqrpc struct {
	*generator.Generator
}

func NewMqRpc() *mqrpc {
	return &mqrpc{}
}

func (p *mqrpc) Name() string {
	return "mqrpc"
}

func (p *mqrpc) Init(g *generator.Generator) {
	p.Generator = g
}

// GenerateImports generates the import declaration for this file.
func (p *mqrpc) GenerateImports(file *generator.FileDescriptor) {
	imports := generator.NewPluginImports(p.Generator)
	imports.NewImport("fmt").Use()
	imports.NewImport("errors").Use()
	imports.NewImport("github.com/gavrilaf/amqp/rpc").Use()
	imports.GenerateImports(file)
}

func (p *mqrpc) Generate(file *generator.FileDescriptor) {

	for _, service := range file.Service {

		p.P(`// Server API`)
		serverInterfaceName := service.GetName() + "Server"
		p.P(`type `, serverInterfaceName, ` interface {`)
		p.In()
		for _, method := range service.GetMethod() {
			methodName := method.GetName()
			argType := p.typeName(method.GetInputType())
			retType := p.typeName(method.GetOutputType())
			p.P(methodName, `(arg *`, argType, `) (*`, retType, `, error)`)
		}
		p.Out()
		p.P(`}`)

		p.P(`// Run server API with this call`)
		p.P(`func RunServer(srv rpc.Server, handler `, serverInterfaceName, `) {`)
		p.In()
		p.P(`srv.Serve(func(funcID int32, arg []byte) ([]byte, error) {`)
		p.In()
		p.P(`switch funcID {`)
		for _, method := range service.GetMethod() {
			methodName := method.GetName()
			p.P(`case Functions_`, methodName, `:`)
			p.In()
			p.P(`return _Handle_`, methodName, `(handler, arg)`)
			p.Out()
		}
		p.P(`default:`)
		p.In()
		p.P(`return nil, errors.New(fmt.Sprintf("unknown function with code: %d", funcID))`)
		p.Out()
		p.P(`}`)
		p.P(`})`)
		p.Out()
		p.P(`}`)

		p.P(`// Client API`)
		clientInterfaceName := service.GetName() + "Client"
		p.P(`type `, clientInterfaceName, ` interface {`)
		p.In()
		for _, method := range service.GetMethod() {
			methodName := method.GetName()
			argType := p.typeName(method.GetInputType())
			retType := p.typeName(method.GetOutputType())
			p.P(methodName, `(arg *`, argType, `) (*`, retType, `, error)`)
		}
		p.Out()
		p.P(`}`)

		p.P(`func New`, clientInterfaceName, `(cc rpc.Client)`, clientInterfaceName, `{`)
		p.In()
		p.P(`return &`, unexport(clientInterfaceName), `{cc}`)
		p.Out()
		p.P(`}`)

		p.P(`type `, unexport(clientInterfaceName), ` struct {`)
		p.In()
		p.P(`cc rpc.Client`)
		p.Out()
		p.P(`}`)

		p.P(`// Functions enum`)
		p.P(`const (`)
		p.In()
		for indx, method := range service.GetMethod() {
			p.P(`Functions_`, method.GetName(), ` int32 = `, indx)
		}
		p.Out()
		p.P(`)`)

		p.P(`// Server API handlers`)
		for _, method := range service.GetMethod() {
			methodName := method.GetName()
			argType := p.typeName(method.GetInputType())

			p.P(`func _Handle_`, methodName, `(handler interface{}, arg []byte) ([]byte, error) {`)
			p.In()
			p.P(`var req `, argType)
			p.P(`err := req.Unmarshal(arg)`)
			p.printErr()

			p.P(`resp, err := handler.(`, serverInterfaceName, `).`, methodName, `(&req)`)
			p.printErr()

			p.P(`return resp.Marshal()`)
			p.Out()
			p.P(`}`)
		}

		p.P(`// Client API handlers`)
		for _, method := range service.GetMethod() {
			methodName := method.GetName()
			argType := p.typeName(method.GetInputType())
			retType := p.typeName(method.GetOutputType())

			p.P(`func (this *`, unexport(clientInterfaceName), `)`, methodName, `(arg *`, argType, `) (*`, retType, `, error) {`)
			p.In()
			p.P(`request, err := arg.Marshal()`)
			p.printErr()
			p.P(`respData, err := this.cc.RemoteCall(rpc.Request{FuncID: Functions_`, methodName, `, Body: request})`)
			p.printErr()

			p.P(`var resp `, retType)
			p.P(`err = resp.Unmarshal(respData)`)
			p.P(`return &resp, err`)
			p.Out()
			p.P(`}`)
		}
	}
}

func (p *mqrpc) typeName(s string) string {
	return p.Generator.TypeName(p.Generator.ObjectNamed(s))
}

func (p *mqrpc) printErr() {
	p.P(`if err != nil {`)
	p.In()
	p.P(`return nil, err`)
	p.Out()
	p.P(`}`)
}

func unexport(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func init() {
	generator.RegisterPlugin(NewMqRpc())
}
