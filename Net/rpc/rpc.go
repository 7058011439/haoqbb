package rpc

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"google.golang.org/grpc"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	CodeTypeGob int = iota
	CodeTypeJson
)

type IRpc interface {
	Name() string
}

func NewGRpcServer(port int) *GServer {
	return &GServer{
		port:   port,
		Server: grpc.NewServer(),
	}
}

func NewRpcServer(port int, codeType int) *Server {
	return &Server{
		port:     port,
		codeType: codeType,
	}
}

func NewRpcClient(addr string, port int, i IRpc, codeType int) *Client {
	var client *rpc.Client
	var err error
	switch codeType {
	case CodeTypeGob:
		if client, err = rpc.Dial("tcp", fmt.Sprintf("%v:%v", addr, port)); err != nil {
			return nil
		}
	case CodeTypeJson:
		if conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", addr, port)); err != nil {
			return nil
		} else {
			client = rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
		}
	}
	return &Client{
		Client: client,
		name:   i.Name(),
	}
}

type Server struct {
	port     int
	codeType int
}

func (r *Server) RegisterService(i IRpc) error {
	return rpc.RegisterName(i.Name(), i)
}

func (r *Server) Start() error {
	if listen, err := net.Listen("tcp", fmt.Sprintf(":%v", r.port)); err != nil {
		return err
	} else {
		go func() {
			for {
				if accept, err := listen.Accept(); err != nil {
					Log.Error("监听连接错误, err = %v", err)
				} else {
					switch r.codeType {
					case CodeTypeGob:
						go rpc.ServeConn(accept)
					case CodeTypeJson:
						go rpc.ServeCodec(jsonrpc.NewServerCodec(accept))
					}
				}
			}
		}()
	}
	return nil
}

type Client struct {
	*rpc.Client
	name string
}

func (r *Client) Call(method string, input interface{}, out interface{}) {
	r.Client.Call(r.name+"."+method, input, out)
}
