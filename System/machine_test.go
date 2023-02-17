package System

import (
	"fmt"
	"testing"
)

func TestGetMachineInfo(t *testing.T) {
	a := GetMachineInfo()
	fmt.Println(a)

	c := GetMachineId()
	fmt.Println(c)
}
