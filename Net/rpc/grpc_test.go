package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"testing"
)

type CalcServer struct {
	UnimplementedCalcServiceServer
}

func (c *CalcServer) Multiply(ctx context.Context, req *CalcRequest) (*CalcResponse, error) {
	result := req.A * req.B
	return &CalcResponse{Result: result}, nil
}

func (c *CalcServer) Divide(ctx context.Context, req *CalcRequest) (*CalcResponse, error) {
	if req.B == 0 {
		return nil, fmt.Errorf("divide by zero")
	}
	result := req.A / req.B
	return &CalcResponse{Result: result}, nil
}

func TestGRpc(t *testing.T) {
	port := 666
	s := NewGRpcServer(port)
	s.RegisterService(&CalcService_ServiceDesc, &CalcServer{})
	s.Start()

	conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%v", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := NewCalcServiceClient(conn)

	requestData := &CalcRequest{
		A: rand.Int31n(10086),
		B: rand.Int31n(1008) + 1,
	}
	// Multiply
	mulResp, err := client.Multiply(context.Background(), requestData)
	if err != nil {
		log.Fatalf("Error during Multiply: %v", err)
	}
	log.Printf("%d * %d = %d", requestData.A, requestData.B, mulResp.Result)

	// Divide
	divResp, err := client.Divide(context.Background(), requestData)
	if err != nil {
		log.Fatalf("Error during Divide: %v", err)
	}
	log.Printf("%v / %v = %v", requestData.A, requestData.B, divResp.Result)

	select {}
}
