package greet

import (
	"gola/grpc/discover"
	"gola/internal/logger"
	"net/rpc"
	"sync"
)

type rpcClient struct {
	sync.Once
	client *rpc.Client
}

func (c *rpcClient) New() *rpc.Client {
	c.Do(func() {
		client, err := rpc.Dial("tcp", discover.Discover("greet"))
		if err != nil {
			panic(err)
		}
		c.client = client
	})
	return c.client
}

var greetClient = &rpcClient{}

func Client(name string) (*HelloReply, error) {
	client := greetClient.New()

	var res = &HelloReply{}
	err := client.Call("GreetServer.SayHello", &HelloRequest{Name: name}, &res)
	if err != nil {
		return nil, err
	}
	logger.Success(res.Message)

	err = client.Call("GreetServer.SayHelloAgain", &HelloRequest{Name: name}, &res)
	if err != nil {
		return nil, err
	}
	logger.Success(res.Message)

	return res, nil
}
