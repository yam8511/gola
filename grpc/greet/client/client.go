package greet

import (
	"context"
	"gola/grpc/discover"
	"gola/grpc/proto/greet"
	"gola/internal/logger"
	"sync"

	"google.golang.org/grpc"
)

type grpcClient struct {
	sync.Once
	conn *grpc.ClientConn
}

func (c *grpcClient) New() *grpc.ClientConn {
	c.Do(func() {
		conn, err := grpc.Dial(discover.Discover("greet", "grpc"), grpc.WithInsecure())
		if err != nil {
			panic(err)
		}
		c.conn = conn
	})
	return c.conn
}

var greetClient = &grpcClient{}

func Client(name string) (*greet.HelloReply, error) {

	client := greet.NewGreeterClient(greetClient.New())

	res, err := client.SayHello(
		context.Background(),
		&greet.HelloRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	logger.Success(res.GetMessage())

	res, err = client.SayHelloAgain(
		context.Background(),
		&greet.HelloRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	logger.Success(res.GetMessage())

	return res, nil
}
