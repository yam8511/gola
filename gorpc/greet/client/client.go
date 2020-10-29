package greet

import (
	"gola/gorpc/proto"
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
		client, err := rpc.Dial("tcp", discover.Discover("greet", "gorpc"))
		if err != nil {
			panic(err)
		}
		c.client = client
	})
	return c.client
}

var greetClient = &rpcClient{}

func Client(name string) (*proto.HelloReply, error) {
	client := greetClient.New()

	var res = &proto.HelloReply{}
	err := client.Call("GreetServer.SayHello", &proto.HelloRequest{Name: name}, &res)
	if err != nil {
		return nil, err
	}
	logger.Success(res.Message)

	err = client.Call("GreetServer.SayHelloAgain", &proto.HelloRequest{Name: name}, &res)
	if err != nil {
		return nil, err
	}
	logger.Success(res.Message)

	return res, nil
}
