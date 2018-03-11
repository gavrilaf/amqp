package test

import (
	"errors"
	"github.com/gavrilaf/amqp/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var errTest = errors.New("testErr")

type testSrv struct{}

func (p testSrv) CopySimple(arg *SimpleTypes) (*SimpleTypes, error) {
	copy := *arg
	return &copy, nil
}

func (p testSrv) GenErr(arg *Empty) (*Empty, error) {
	return nil, errTest
}

func runSrv(t *testing.T) *rpc.Server {
	srv, err := rpc.CreateServer("amqp://localhost:5672", "rpc-test-worker")
	require.Nil(t, err, "Couldn't create server")

	go func(s *rpc.Server) {
		RunServer(s, &testSrv{})
	}(srv)

	return srv
}

func clientConnect(t *testing.T) (TestClient, rpc.Client) {
	cc, err := rpc.Connect(rpc.ClientConfig{Url: "amqp://localhost:5672", ServerQueue: "rpc-test-worker", Timeout: time.Second})
	require.Nil(t, err, "Couldn't connect client")

	return NewTestClient(cc), cc
}

func Test_CopySimple(t *testing.T) {
	runSrv(t)
	//defer srv.Close()

	client, cc := clientConnect(t)
	defer cc.Close()

	arg := SimpleTypes{Number: 12, Str: "String", Logic: true}
	res, err := client.CopySimple(&arg)

	assert.Nil(t, err)
	assert.True(t, arg.Equal(res))
}

func Test_Error(t *testing.T) {
	runSrv(t)
	//defer srv.Close()

	client, cc := clientConnect(t)
	defer cc.Close()

	res, err := client.GenErr(&Empty{})
	assert.Nil(t, res)
	assert.NotNil(t, err)

	assert.Equal(t, err.Error(), errTest.Error())
}
