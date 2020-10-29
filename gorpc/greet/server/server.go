package greet

import (
	"gola/global"
	"gola/gorpc/proto"
	"gola/internal/bootstrap"
	"gola/internal/logger"
	"gola/internal/server"
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type GreetServer struct {
	mx    *sync.RWMutex
	count int
}

func (gs *GreetServer) SayHello(req *proto.HelloRequest, res *proto.HelloReply) error {
	gs.mx.Lock()
	gs.count++
	num := strconv.Itoa(gs.count)
	gs.mx.Unlock()

	msg := "#" + num + ": Get Hello " + global.AppVersion + " from " + req.Name
	logger.Success(msg)

	res.Message = msg
	return nil
}

func (gs *GreetServer) SayHelloAgain(req *proto.HelloRequest, res *proto.HelloReply) error {
	gs.mx.Lock()
	gs.count++
	num := strconv.Itoa(gs.count)
	gs.mx.Unlock()

	msg := "#" + num + ": Get Hello Again " + global.AppVersion + " from " + req.Name
	logger.Success(msg)

	// time.Sleep(time.Second * 2)

	res.Message = msg
	return nil
	// return nil, status.Errorf(codes.Unimplemented, "method SayHelloAgain not implemented")
}

func Server() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	dl := server.NewDozListener(l, 0, true)

	server := rpc.NewServer()
	service := &GreetServer{&sync.RWMutex{}, 0}
	err = server.Register(service)
	if err != nil {
		panic(err)
	}

	logger.Info("Greet Service Listening... " + l.Addr().String())

	go func() {
		for {
			conn, err := dl.Accept()
			if err != nil {
				log.Println("Greet Service Stop: ", err)
				return
			}

			go server.ServeConn(conn)
		}
	}()

	<-bootstrap.GracefulDown()
	err = dl.Close()
	if err != nil {
		logger.Danger("DozListener Close 發生錯誤: %s", err.Error())
	}
	dl.Wait()
	if err != nil {
		logger.Danger("DozListener Wait 發生錯誤: %s", err.Error())
	}

	// ctx, done := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}
	// go func() {
	// 	defer func() {
	// 		wg.Wait()
	// 		done()
	// 	}()

	// 	for {
	// 		conn, err := l.Accept()
	// 		if err != nil {
	// 			log.Println("Greet Service Stop: ", err)
	// 			return
	// 		}

	// 		wg.Add(1)
	// 		go func() {
	// 			server.ServeConn(conn)
	// 			wg.Done()
	// 		}()
	// 	}
	// }()

	// <-bootstrap.GracefulDown()
	// l.Close()
	// <-ctx.Done()
}
