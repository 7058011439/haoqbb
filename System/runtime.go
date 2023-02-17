package System

import (
	"Core/Log"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func DropErr(e error, funName string) {
	if e != nil {
		Log.ErrorLog("[%s]:%s", funName, e)
		panic(e)
	}
}

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func StartExe(exeName string, params ...string) {
	cmd := exec.Command("cmd.exe", "/C", "start", exeName)
	if err := cmd.Run(); err != nil {
		Log.ErrorLog("启动进程失败, err = %v", err)
	}
}

func GetCurrentDirectory() string {
	//返回绝对路径 filepath.Dir(os.Args[0])去除最后一个元素的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	//将\替换成/
	return strings.Replace(dir, "\\", "/", -1)
}
