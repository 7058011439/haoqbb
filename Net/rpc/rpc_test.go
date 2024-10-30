package rpc

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type helloService struct{}

func (h *helloService) Name() string {
	return "hello"
}

func (h *helloService) SayHi(request string, response *string) error {
	*response = fmt.Sprintf("hi %v %v", request, time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

// 定义输入参数
type Args struct {
	A, B int
}

type calc struct{}

func (t *calc) Name() string {
	return "calc"
}

func (t *calc) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *calc) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = args.A / args.B
	return nil
}

func (t *calc) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func (t *calc) Sub(args *Args, reply *int) error {
	*reply = args.A - args.B
	return nil
}

func TestRpc(t *testing.T) {
	port := 6666
	s := NewRpcServer(port, CodeTypeJson)
	s.RegisterService(&helloService{})
	s.RegisterService(&calc{})
	s.Start()

	cHello := NewRpcClient("127.0.0.1", port, &helloService{}, CodeTypeJson)
	var retStr string
	for i := 0; i < 3; i++ {
		cHello.Call("SayHi", "superman", &retStr)
		fmt.Println(retStr)
		time.Sleep(time.Second)
	}

	cCalc := NewRpcClient("127.0.0.1", port, &calc{}, CodeTypeJson)
	var retInt int
	for i := 0; i < 2; i++ {
		data := &Args{rand.Intn(10086), rand.Intn(10086)}
		cCalc.Call("Multiply", data, &retInt)
		fmt.Printf("%v * %v = %v\n", data.A, data.B, retInt)
		cCalc.Call("Divide", data, &retInt)
		fmt.Printf("%v / %v = %v\n", data.A, data.B, retInt)
		cCalc.Call("Add", data, &retInt)
		fmt.Printf("%v + %v = %v\n", data.A, data.B, retInt)
		cCalc.Call("Sub", data, &retInt)
		fmt.Printf("%v - %v = %v\n", data.A, data.B, retInt)
		fmt.Println()
	}
}
