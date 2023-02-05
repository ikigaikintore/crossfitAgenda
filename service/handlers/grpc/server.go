package grpc

import (
	grpcLib "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"log"
	"net"
	"sync"
	"time"
)

type CrossfitServer struct {
	ls  net.Listener
	srv *grpcLib.Server
}

func NewServer(handler CrossfitAgendaSvc) *CrossfitServer {
	ls, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	srv := grpcLib.NewServer(
		grpcLib.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 25 * time.Second,
			MaxConnectionAge:  41 * time.Second,
			Time:              20 * time.Minute,
			Timeout:           15 * time.Second,
		}),
		grpcLib.ConnectionTimeout(50*time.Second),
	)

	RegisterDefaultServiceServer(srv, handler)

	return &CrossfitServer{
		ls:  ls,
		srv: srv,
	}
}

func (cs *CrossfitServer) Serve() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cs.srv.Serve(cs.ls); err != nil {
			log.Fatalln(err)
		}
	}()
	wg.Wait()
}

func (cs *CrossfitServer) Shutdown() {
	cs.srv.GracefulStop()
}
