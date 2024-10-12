package protocol

import (
	"github.com/7058011439/haoqbb/String"
	protoFast "github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/proto"
	"math/rand"
	"testing"
)

var testLoginData *C2S_LoginWithToken

func initTestData() {
	testLoginData = &C2S_LoginWithToken{
		MachineId:     String.RandStr(30),
		Token:         String.RandStr(20),
		Phone:         String.RandStr(11),
		SrvId:         rand.Int31n(100),
		Channel:       rand.Int31n(100),
		GameId:        rand.Int31n(100),
		MainVer:       rand.Int31n(100),
		EvaluationVer: rand.Int31n(100),
		HotfixVer:     rand.Int31n(100),
	}
}

/*
当前基准测试显示使用gogo(快速)的效率比golang(标准)提高约200%
BenchmarkC2S_LoginWithTokenProto
BenchmarkC2S_LoginWithTokenProto-16     23907115               505.6 ns/op
BenchmarkC2S_LoginWithTokenGogo
BenchmarkC2S_LoginWithTokenGogo-16      64569728               178.3 ns/op
*/
func BenchmarkC2S_LoginWithTokenProto(b *testing.B) {
	initTestData()
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = proto.Marshal(testLoginData)
		proto.Unmarshal(data, testLoginData)
	}
}

func BenchmarkC2S_LoginWithTokenGogo(b *testing.B) {
	initTestData()
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = protoFast.Marshal(testLoginData)
		protoFast.Unmarshal(data, testLoginData)
	}
}
