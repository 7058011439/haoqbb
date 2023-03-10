package Util

import (
	cr "crypto/rand"
	"fmt"
	"io"
)

//func GetUid() string {
//	b := make([]byte, 16)
//	io.ReadFull(cr.Reader, b)
//	b[6] = (b[6] & 0x0f) | 0x40
//	b[8] = (b[8] & 0x3f) | 0x80
//	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
//}

func GetUid() string {
	b := make([]byte, 16)
	io.ReadFull(cr.Reader, b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x", b)
}
