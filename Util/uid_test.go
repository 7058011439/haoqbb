package Util

import "testing"

func TestGetUid(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(GetUid())
	}
}
