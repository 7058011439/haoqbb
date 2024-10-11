package Panic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRedirectStderr(t *testing.T) {
	RedirectStderr("./Logs/Panic.log")
	go func() {
		arrInt := []int{0, 1}
		index := 0
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-tick.C:
				fmt.Println("index", arrInt[index])
				index++
			}
		}
	}()
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-tick.C:
				fmt.Println("div", rand.Intn(10)/rand.Intn(5))
			}
		}
	}()
	select {}
	fmt.Println(111)
}
