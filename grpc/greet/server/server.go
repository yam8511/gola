package greet

import (
	"context"
	"gola/global"
	"gola/grpc/proto/greet"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"gola/internal/server"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type GreetServer struct {
	greet.UnimplementedGreeterServer
	mx    *sync.RWMutex
	count int
}

func (*GreetServer) SayHello(ctx context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	msg := "Get Hello " + global.AppVersion + " from " + req.String()
	logger.Success(msg)
	res := &greet.HelloReply{
		Message: msg,
	}
	return res, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}

func (*GreetServer) SayHelloAgain(_ context.Context, req *greet.HelloRequest) (*greet.HelloReply, error) {
	msg := "Get Hello Again " + global.AppVersion + " from " + req.String()
	logger.Success(msg)

	// time.Sleep(time.Second * 2)

	res := &greet.HelloReply{
		Message: msg,
	}

	return res, nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHelloAgain not implemented")
}

func Server() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	dl := server.NewDozListener(l, 0, true)

	server := grpc.NewServer()
	service := &GreetServer{mx: &sync.RWMutex{}}
	greet.RegisterGreeterServer(server, service)

	logger.Info("Greet Service Listening... " + l.Addr().String())

	go func() {
		// err := server.Serve(l)
		err := server.Serve(dl)
		log.Println("Greet Service Stop: ", err)
	}()

	<-bootstrap.GracefulDown()
	server.GracefulStop()
	err = dl.Wait()
	if err != nil {
		logger.Danger("Wait失敗: %s", err.Error())
	}
}
