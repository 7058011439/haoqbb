package Stl

import (
	"github.com/7058011439/haoqbb/Util"
	"github.com/pkg/errors"
)

func NewBuffer(size int) *Buffer {
	ret := &Buffer{
		bs: make([]byte, 0, size),
	}
	ret.cs = ret.bs
	return ret
}

type Buffer struct {
	bs []byte
	cs []byte
}

func (b *Buffer) Write(data []byte) {
	b.checkCap(len(data))
	b.cs = append(b.cs, data...)
}

func (b *Buffer) WriteInt16(data int16) {
	b.Write(Util.Int16ToBytes(data))
}

func (b *Buffer) WriteUint16(data uint16) {
	b.Write(Util.Uint16ToBytes(data))
}

func (b *Buffer) WriteInt(data int) {
	b.Write(Util.IntToBytes(data))
}

func (b *Buffer) WriteUInt(data uint) {
	b.Write(Util.UintToBytes(data))
}

func (b *Buffer) WriteInt32(data int32) {
	b.Write(Util.Int32ToBytes(data))
}

func (b *Buffer) WriteUInt32(data uint32) {
	b.Write(Util.Uint32ToBytes(data))
}

func (b *Buffer) WriteInt64(data int64) {
	b.Write(Util.Int64ToBytes(data))
}

func (b *Buffer) WriteUInt64(data uint64) {
	b.Write(Util.Uint64ToBytes(data))
}

func (b *Buffer) WriteFloat32(data float32) {
	b.Write(Util.Float32ToBytes(data))
}

func (b *Buffer) WriteFloat64(data float64) {
	b.Write(Util.Float64ToBytes(data))
}

func (b *Buffer) WriteByte(data byte) {
	b.checkCap(1)
	b.cs = append(b.cs, data)
}

func (b *Buffer) WriteString(data string) {
	b.Write([]byte(data))
}

func (b *Buffer) Bytes() []byte {
	return b.cs
}

func (b *Buffer) String() string {
	return string(b.cs)
}

func (b *Buffer) Len() int {
	return len(b.cs)
}

func (b *Buffer) Cap() int {
	return cap(b.cs)
}

func (b *Buffer) OffSize(os int) error {
	if os > len(b.cs) {
		return errors.Errorf("Failed to OffSize, os larger than size")
	}
	buff := b.cs[os:]
	// 这里很绕，既要让cs 不重新分配内存空间，又要让cs有足够容量存放保留内容
	b.cs = b.bs[0 : len(b.cs)-os]
	copy(b.cs, buff)
	return nil
}

func (b *Buffer) Reset() {
	b.cs = b.bs
}

func (b *Buffer) checkCap(size int) {
	// 这个地方写成 size+len(b.cs)>cap(b.cs) 也可以
	if size+len(b.cs) > cap(b.bs) {
		b.bs = make([]byte, 0, (size+len(b.cs))*2)
		b.OffSize(0)
	}
}
