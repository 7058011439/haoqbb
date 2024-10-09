package Stl

import (
	"github.com/7058011439/haoqbb/Util"
	"github.com/pkg/errors"
)

func NewBuffer(size int) *Buffer {
	ret := &Buffer{
		cs: make([]byte, 0, size),
	}
	return ret
}

type Buffer struct {
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
	b.cs = b.cs[:copy(b.cs, b.cs[os:])]
	return nil
}

func (b *Buffer) Reset() {
	b.cs = b.cs[:0]
}

func (b *Buffer) checkCap(size int) {
	required := size + len(b.cs)
	if required > cap(b.cs) {
		if required < 1024 {
			required *= 2
		} else {
			required = int(float64(required) * 1.2)
		}
		// 直接在扩容后的新切片中赋值
		b.cs = append(make([]byte, 0, required), b.cs...)
	}
}
