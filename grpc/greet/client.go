package greet

import (
	"context"
	"gola/grpc/discover"
	"gola/grpc/proto/greet"
	"gola/internal/logger"

	"google.golang.org/grpc"
)

func Client(name string) (*greet.HelloReply, error) {
	conn, err := grpc.Dial(discover.Discover("greet", "grpc"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := greet.NewGreeterClient(conn)

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
