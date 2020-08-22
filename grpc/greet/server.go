package greet

import (
	"context"
	"gola/grpc/proto/greet"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GreetServer struct {
	greet.UnimplementedGreeterServer
}

func (*GreetServer) SayHello(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	logger.Success("Get Hello from " + req.String())
	res := &greet.HelloReply{
		Message: "Hello " + req.GetName(),
	}
	return res, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

func (*GreetServer) SayHelloAgain(_ context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	logger.Success("Get Hello Again from " + req.String())

	time.Sleep(time.Second * 2)

	res := &greet.HelloReply{
		Message: "Hello Again " + req.GetName(),
	}

	return res, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHelloAgain not implemented")
}

func Server() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	service := new(GreetServer)
	greet.RegisterGreeterServer(server, service)

	logger.Info("Greet Service Listening... " + l.Addr().String())

	ctx, done := context.WithCancel(context.Background())
	go func() {
		err := server.Serve(l)
		log.Println("Greet Service Stop: ", err)
		done()
	}()

	<-bootstrap.GracefulDown()
	server.GracefulStop()
	<-ctx.Done()
}
