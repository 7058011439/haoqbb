package System

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
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
}

func TestRedirectStderr(t *testing.T) {
	RedirectStderr("./Logs/panic.log")
	go func() {
		tick := time.NewTicker(time.Second)
		for {
			select {
			case <-tick.C:
				fmt.Println("div", rand.Intn(10)/rand.Intn(5))
			}
		}
	}()
	chStop := make(chan os.Signal)
	<-chStop
	fmt.Println(111)
}
