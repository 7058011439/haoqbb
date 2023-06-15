// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: protoGameServer.proto

package protocol

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type C2S_Test_RT struct {
	Index                int64    `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *C2S_Test_RT) Reset()         { *m = C2S_Test_RT{} }
func (m *C2S_Test_RT) String() string { return proto.CompactTextString(m) }
func (*C2S_Test_RT) ProtoMessage()    {}
func (*C2S_Test_RT) Descriptor() ([]byte, []int) {
	return fileDescriptor_88cc3d0b27d75a2e, []int{0}
}
func (m *C2S_Test_RT) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *C2S_Test_RT) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_C2S_Test_RT.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *C2S_Test_RT) XXX_Merge(src proto.Message) {
	xxx_messageInfo_C2S_Test_RT.Merge(m, src)
}
func (m *C2S_Test_RT) XXX_Size() int {
	return m.Size()
}
func (m *C2S_Test_RT) XXX_DiscardUnknown() {
	xxx_messageInfo_C2S_Test_RT.DiscardUnknown(m)
}

var xxx_messageInfo_C2S_Test_RT proto.InternalMessageInfo

func (m *C2S_Test_RT) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

type S2C_Test_RT struct {
	Index                int64    `protobuf:"varint,1,opt,name=index,proto3" json:"index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *S2C_Test_RT) Reset()         { *m = S2C_Test_RT{} }
func (m *S2C_Test_RT) String() string { return proto.CompactTextString(m) }
func (*S2C_Test_RT) ProtoMessage()    {}
func (*S2C_Test_RT) Descriptor() ([]byte, []int) {
	return fileDescriptor_88cc3d0b27d75a2e, []int{1}
}
func (m *S2C_Test_RT) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *S2C_Test_RT) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_S2C_Test_RT.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *S2C_Test_RT) XXX_Merge(src proto.Message) {
	xxx_messageInfo_S2C_Test_RT.Merge(m, src)
}
func (m *S2C_Test_RT) XXX_Size() int {
	return m.Size()
}
func (m *S2C_Test_RT) XXX_DiscardUnknown() {
	xxx_messageInfo_S2C_Test_RT.DiscardUnknown(m)
}

var xxx_messageInfo_S2C_Test_RT proto.InternalMessageInfo

func (m *S2C_Test_RT) GetIndex() int64 {
	if m != nil {
		return m.Index
	}
	return 0
}

//通过TOKEN登陆游戏服务器,machineId,token,srvid必须要
type C2S_LoginWithToken struct {
	MachineId            string   `protobuf:"bytes,1,opt,name=machineId,proto3" json:"machineId,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	SrvId                int32    `protobuf:"varint,3,opt,name=srvId,proto3" json:"srvId,omitempty"`
	Channel              int32    `protobuf:"varint,4,opt,name=channel,proto3" json:"channel,omitempty"`
	GameId               int32    `protobuf:"varint,5,opt,name=gameId,proto3" json:"gameId,omitempty"`
	MainVer              int32    `protobuf:"varint,6,opt,name=mainVer,proto3" json:"mainVer,omitempty"`
	EvaluationVer        int32    `protobuf:"varint,7,opt,name=evaluationVer,proto3" json:"evaluationVer,omitempty"`
	HotfixVer            int32    `protobuf:"varint,8,opt,name=hotfixVer,proto3" json:"hotfixVer,omitempty"`
	Phone                string   `protobuf:"bytes,9,opt,name=phone,proto3" json:"phone,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *C2S_LoginWithToken) Reset()         { *m = C2S_LoginWithToken{} }
func (m *C2S_LoginWithToken) String() string { return proto.CompactTextString(m) }
func (*C2S_LoginWithToken) ProtoMessage()    {}
func (*C2S_LoginWithToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_88cc3d0b27d75a2e, []int{2}
}
func (m *C2S_LoginWithToken) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *C2S_LoginWithToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_C2S_LoginWithToken.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *C2S_LoginWithToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_C2S_LoginWithToken.Merge(m, src)
}
func (m *C2S_LoginWithToken) XXX_Size() int {
	return m.Size()
}
func (m *C2S_LoginWithToken) XXX_DiscardUnknown() {
	xxx_messageInfo_C2S_LoginWithToken.DiscardUnknown(m)
}

var xxx_messageInfo_C2S_LoginWithToken proto.InternalMessageInfo

func (m *C2S_LoginWithToken) GetMachineId() string {
	if m != nil {
		return m.MachineId
	}
	return ""
}

func (m *C2S_LoginWithToken) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *C2S_LoginWithToken) GetSrvId() int32 {
	if m != nil {
		return m.SrvId
	}
	return 0
}

func (m *C2S_LoginWithToken) GetChannel() int32 {
	if m != nil {
		return m.Channel
	}
	return 0
}

func (m *C2S_LoginWithToken) GetGameId() int32 {
	if m != nil {
		return m.GameId
	}
	return 0
}

func (m *C2S_LoginWithToken) GetMainVer() int32 {
	if m != nil {
		return m.MainVer
	}
	return 0
}

func (m *C2S_LoginWithToken) GetEvaluationVer() int32 {
	if m != nil {
		return m.EvaluationVer
	}
	return 0
}

func (m *C2S_LoginWithToken) GetHotfixVer() int32 {
	if m != nil {
		return m.HotfixVer
	}
	return 0
}

func (m *C2S_LoginWithToken) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

// 登录成功
type S2C_GameLoginResult struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Err                  string   `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
	Code                 int32    `protobuf:"varint,4,opt,name=code,proto3" json:"code,omitempty"`
	ServerTimeNow        uint64   `protobuf:"varint,5,opt,name=serverTimeNow,proto3" json:"serverTimeNow,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *S2C_GameLoginResult) Reset()         { *m = S2C_GameLoginResult{} }
func (m *S2C_GameLoginResult) String() string { return proto.CompactTextString(m) }
func (*S2C_GameLoginResult) ProtoMessage()    {}
func (*S2C_GameLoginResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_88cc3d0b27d75a2e, []int{3}
}
func (m *S2C_GameLoginResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *S2C_GameLoginResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_S2C_GameLoginResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *S2C_GameLoginResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_S2C_GameLoginResult.Merge(m, src)
}
func (m *S2C_GameLoginResult) XXX_Size() int {
	return m.Size()
}
func (m *S2C_GameLoginResult) XXX_DiscardUnknown() {
	xxx_messageInfo_S2C_GameLoginResult.DiscardUnknown(m)
}

var xxx_messageInfo_S2C_GameLoginResult proto.InternalMessageInfo

func (m *S2C_GameLoginResult) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *S2C_GameLoginResult) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func (m *S2C_GameLoginResult) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *S2C_GameLoginResult) GetServerTimeNow() uint64 {
	if m != nil {
		return m.ServerTimeNow
	}
	return 0
}

type Message struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_88cc3d0b27d75a2e, []int{4}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Message.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*C2S_Test_RT)(nil), "protocol.C2S_Test_RT")
	proto.RegisterType((*S2C_Test_RT)(nil), "protocol.S2C_Test_RT")
	proto.RegisterType((*C2S_LoginWithToken)(nil), "protocol.C2S_LoginWithToken")
	proto.RegisterType((*S2C_GameLoginResult)(nil), "protocol.S2C_GameLoginResult")
	proto.RegisterType((*Message)(nil), "protocol.Message")
}

func init() { proto.RegisterFile("protoGameServer.proto", fileDescriptor_88cc3d0b27d75a2e) }

var fileDescriptor_88cc3d0b27d75a2e = []byte{
	// 351 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0xfd, 0xd2, 0xff, 0xcc, 0x87, 0x50, 0xc6, 0x1f, 0x02, 0x4a, 0x28, 0xd5, 0x45, 0x57, 0x2e,
	0xea, 0x1b, 0xd8, 0x85, 0x14, 0xd4, 0xc5, 0x34, 0xe8, 0xb2, 0x8c, 0xc9, 0x35, 0x19, 0x4c, 0x66,
	0x4a, 0x66, 0x5a, 0x0b, 0x6e, 0x7d, 0x08, 0x1f, 0xc9, 0xa5, 0x8f, 0x20, 0xf5, 0x45, 0xe4, 0xde,
	0xa4, 0x94, 0xae, 0x5c, 0xe5, 0x9e, 0x73, 0xcf, 0xbd, 0x39, 0xf7, 0x24, 0xec, 0x78, 0x51, 0x1a,
	0x67, 0x6e, 0x64, 0x01, 0x33, 0x28, 0x57, 0x50, 0x5e, 0x12, 0xe6, 0x3d, 0x7a, 0xc4, 0x26, 0x1f,
	0x9e, 0xb3, 0xff, 0x93, 0xf1, 0x6c, 0x1e, 0x81, 0x75, 0x73, 0x11, 0xf1, 0x23, 0xd6, 0x56, 0x3a,
	0x81, 0x75, 0xe0, 0x0d, 0xbc, 0x51, 0x53, 0x54, 0x00, 0x45, 0xb3, 0xf1, 0xe4, 0x0f, 0xd1, 0x7b,
	0x83, 0x71, 0x5c, 0x75, 0x6b, 0x52, 0xa5, 0x1f, 0x95, 0xcb, 0x22, 0xf3, 0x02, 0x9a, 0x9f, 0x31,
	0xbf, 0x90, 0x71, 0xa6, 0x34, 0x4c, 0x13, 0x1a, 0xf0, 0xc5, 0x8e, 0xc0, 0x55, 0x0e, 0x65, 0x41,
	0x83, 0x3a, 0x15, 0x40, 0xd6, 0x96, 0xab, 0x69, 0x12, 0x34, 0x07, 0xde, 0xa8, 0x2d, 0x2a, 0xc0,
	0x03, 0xd6, 0x8d, 0x33, 0xa9, 0x35, 0xe4, 0x41, 0x8b, 0xf8, 0x2d, 0xe4, 0x27, 0xac, 0x93, 0xca,
	0x02, 0x5f, 0xd0, 0xa6, 0x46, 0x8d, 0x70, 0xa2, 0x90, 0x4a, 0x3f, 0x40, 0x19, 0x74, 0xaa, 0x89,
	0x1a, 0xf2, 0x0b, 0x76, 0x00, 0x2b, 0x99, 0x2f, 0xa5, 0x53, 0x86, 0xfa, 0x5d, 0xea, 0xef, 0x93,
	0xe8, 0x3d, 0x33, 0xee, 0x59, 0xad, 0x51, 0xd1, 0x23, 0xc5, 0x8e, 0x40, 0x97, 0x8b, 0xcc, 0x68,
	0x08, 0xfc, 0xca, 0x3b, 0x81, 0xe1, 0x1b, 0x3b, 0xc4, 0xac, 0x30, 0x72, 0x4a, 0x42, 0x80, 0x5d,
	0xe6, 0x0e, 0xad, 0xd8, 0x65, 0x1c, 0x83, 0xb5, 0x14, 0x42, 0x4f, 0x6c, 0x21, 0xef, 0xb3, 0x26,
	0x94, 0x25, 0x9d, 0xea, 0x0b, 0x2c, 0x39, 0x67, 0xad, 0xd8, 0x24, 0x50, 0x5f, 0x49, 0x35, 0x1a,
	0xb6, 0xf4, 0x05, 0x23, 0x55, 0xc0, 0xbd, 0x79, 0xa5, 0x4b, 0x5b, 0x62, 0x9f, 0x1c, 0x9e, 0xb2,
	0xee, 0x1d, 0x58, 0x2b, 0x53, 0xc0, 0xb5, 0x85, 0x4d, 0xeb, 0xc4, 0xb1, 0xbc, 0xee, 0x7f, 0x6e,
	0x42, 0xef, 0x6b, 0x13, 0x7a, 0xdf, 0x9b, 0xd0, 0xfb, 0xf8, 0x09, 0xff, 0x3d, 0x75, 0xe8, 0x37,
	0xb8, 0xfa, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x27, 0x9d, 0x3e, 0xa4, 0x26, 0x02, 0x00, 0x00,
}

func (m *C2S_Test_RT) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C2S_Test_RT) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *C2S_Test_RT) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Index != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *S2C_Test_RT) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2C_Test_RT) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2C_Test_RT) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Index != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.Index))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *C2S_LoginWithToken) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *C2S_LoginWithToken) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *C2S_LoginWithToken) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Phone) > 0 {
		i -= len(m.Phone)
		copy(dAtA[i:], m.Phone)
		i = encodeVarintProtoGameServer(dAtA, i, uint64(len(m.Phone)))
		i--
		dAtA[i] = 0x4a
	}
	if m.HotfixVer != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.HotfixVer))
		i--
		dAtA[i] = 0x40
	}
	if m.EvaluationVer != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.EvaluationVer))
		i--
		dAtA[i] = 0x38
	}
	if m.MainVer != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.MainVer))
		i--
		dAtA[i] = 0x30
	}
	if m.GameId != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.GameId))
		i--
		dAtA[i] = 0x28
	}
	if m.Channel != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.Channel))
		i--
		dAtA[i] = 0x20
	}
	if m.SrvId != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.SrvId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Token) > 0 {
		i -= len(m.Token)
		copy(dAtA[i:], m.Token)
		i = encodeVarintProtoGameServer(dAtA, i, uint64(len(m.Token)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.MachineId) > 0 {
		i -= len(m.MachineId)
		copy(dAtA[i:], m.MachineId)
		i = encodeVarintProtoGameServer(dAtA, i, uint64(len(m.MachineId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *S2C_GameLoginResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *S2C_GameLoginResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *S2C_GameLoginResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ServerTimeNow != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.ServerTimeNow))
		i--
		dAtA[i] = 0x28
	}
	if m.Code != 0 {
		i = encodeVarintProtoGameServer(dAtA, i, uint64(m.Code))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Err) > 0 {
		i -= len(m.Err)
		copy(dAtA[i:], m.Err)
		i = encodeVarintProtoGameServer(dAtA, i, uint64(len(m.Err)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Success {
		i--
		if m.Success {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintProtoGameServer(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProtoGameServer(dAtA []byte, offset int, v uint64) int {
	offset -= sovProtoGameServer(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *C2S_Test_RT) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovProtoGameServer(uint64(m.Index))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *S2C_Test_RT) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Index != 0 {
		n += 1 + sovProtoGameServer(uint64(m.Index))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *C2S_LoginWithToken) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.MachineId)
	if l > 0 {
		n += 1 + l + sovProtoGameServer(uint64(l))
	}
	l = len(m.Token)
	if l > 0 {
		n += 1 + l + sovProtoGameServer(uint64(l))
	}
	if m.SrvId != 0 {
		n += 1 + sovProtoGameServer(uint64(m.SrvId))
	}
	if m.Channel != 0 {
		n += 1 + sovProtoGameServer(uint64(m.Channel))
	}
	if m.GameId != 0 {
		n += 1 + sovProtoGameServer(uint64(m.GameId))
	}
	if m.MainVer != 0 {
		n += 1 + sovProtoGameServer(uint64(m.MainVer))
	}
	if m.EvaluationVer != 0 {
		n += 1 + sovProtoGameServer(uint64(m.EvaluationVer))
	}
	if m.HotfixVer != 0 {
		n += 1 + sovProtoGameServer(uint64(m.HotfixVer))
	}
	l = len(m.Phone)
	if l > 0 {
		n += 1 + l + sovProtoGameServer(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *S2C_GameLoginResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Success {
		n += 2
	}
	l = len(m.Err)
	if l > 0 {
		n += 1 + l + sovProtoGameServer(uint64(l))
	}
	if m.Code != 0 {
		n += 1 + sovProtoGameServer(uint64(m.Code))
	}
	if m.ServerTimeNow != 0 {
		n += 1 + sovProtoGameServer(uint64(m.ServerTimeNow))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovProtoGameServer(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovProtoGameServer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProtoGameServer(x uint64) (n int) {
	return sovProtoGameServer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *C2S_Test_RT) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: C2S_Test_RT: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C2S_Test_RT: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipProtoGameServer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *S2C_Test_RT) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: S2C_Test_RT: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: S2C_Test_RT: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Index", wireType)
			}
			m.Index = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Index |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipProtoGameServer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *C2S_LoginWithToken) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: C2S_LoginWithToken: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: C2S_LoginWithToken: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MachineId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MachineId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Token = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SrvId", wireType)
			}
			m.SrvId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SrvId |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Channel", wireType)
			}
			m.Channel = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Channel |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GameId", wireType)
			}
			m.GameId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GameId |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MainVer", wireType)
			}
			m.MainVer = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MainVer |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EvaluationVer", wireType)
			}
			m.EvaluationVer = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EvaluationVer |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HotfixVer", wireType)
			}
			m.HotfixVer = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HotfixVer |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Phone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Phone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProtoGameServer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *S2C_GameLoginResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: S2C_GameLoginResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: S2C_GameLoginResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Success", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Success = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Err", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Err = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Code", wireType)
			}
			m.Code = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Code |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServerTimeNow", wireType)
			}
			m.ServerTimeNow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ServerTimeNow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipProtoGameServer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProtoGameServer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProtoGameServer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipProtoGameServer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProtoGameServer
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowProtoGameServer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthProtoGameServer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProtoGameServer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProtoGameServer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProtoGameServer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProtoGameServer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProtoGameServer = fmt.Errorf("proto: unexpected end of group")
)
