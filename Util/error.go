package Util

import (
	"fmt"
	"runtime"
	"strconv"
)

func Check(e error, msg string) {
	if e != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			msg = "[" + file + " 行：" + strconv.Itoa(line) + "]：" + msg
			fmt.Println(msg)
		}
		panic(e)
	}
}
