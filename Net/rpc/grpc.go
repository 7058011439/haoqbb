package rpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type GServer struct {
	port int
	*grpc.Server
}

func (g *GServer) Start() error {
	if lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", g.port)); err != nil {
		return err
	} else {
		go g.Server.Serve(lis)
		return nil
	}
}
