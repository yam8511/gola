package greet

import (
	"context"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
	"time"
)

type GreetServer struct {
	mx    *sync.RWMutex
	count int
}

func (gs *GreetServer) SayHello(req *HelloRequest, res *HelloReply) error {
	gs.mx.Lock()
	gs.count++
	logger.Success("#" + strconv.Itoa(gs.count) + ": Get Hello from " + req.Name)
	gs.mx.Unlock()

	res.Message = "Hello " + req.Name
	return nil
}

func (gs *GreetServer) SayHelloAgain(req *HelloRequest, res *HelloReply) error {
	gs.mx.Lock()
	gs.count++
	logger.Success("#" + strconv.Itoa(gs.count) + ": Get Hello Again from " + req.Name)
	gs.mx.Unlock()

	time.Sleep(time.Second * 2)

	res.Message = "Hello Again " + req.Name
	return nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHelloAgain not implemented")
}

func Server() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	server := rpc.NewServer()
	service := &GreetServer{&sync.RWMutex{}, 0}
	err = server.Register(service)
	if err != nil {
		panic(err)
	}

	logger.Info("Greet Service Listening... " + l.Addr().String())

	ctx, done := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	go func() {
		defer func() {
			wg.Wait()
			done()
		}()

		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println("Greet Service Stop: ", err)
				return
			}

			wg.Add(1)
			go func() {
				server.ServeConn(conn)
				wg.Done()
			}()
		}
	}()

	<-bootstrap.GracefulDown()
	l.Close()
	<-ctx.Done()
}
