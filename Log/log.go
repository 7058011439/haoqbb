package Log

import (
	"Core/Stl"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
)

//--------------------------------------------------------------------------------------------------
// 错误日志处理
//--------------------------------------------------------------------------------------------------
var printLevel int
var showFileLine bool
var maxFileSize int64
var boolStarted bool
var dir string

type logLevel = int

const (
	LevelDebug   logLevel = 0
	LevelLog     logLevel = 1
	LevelWarning logLevel = 2
	LevelError   logLevel = 3
)

type logInfo struct {
	desc  string
	color color.Attribute
	queue *Stl.Queue
}

var mapLogInfo = map[int]*logInfo{
	LevelDebug:   {desc: "Debug", color: 34, queue: Stl.NewQueue()},
	LevelLog:     {desc: "Log", color: 32, queue: Stl.NewQueue()},
	LevelWarning: {desc: "Warning", color: 33, queue: Stl.NewQueue()},
	LevelError:   {desc: "Error", color: 31, queue: Stl.NewQueue()},
}

func init() {
	Init(LevelLog, false, 0, "")
}

func Init(iPrintLevel int, bShowFileLine bool, iMaxFileSize int64, logDir string) {
	printLevel = iPrintLevel
	showFileLine = bShowFileLine
	dir = logDir

	if iMaxFileSize < 1 {
		maxFileSize = 5 * 1024 * 1024
	} else {
		maxFileSize = iMaxFileSize
	}
	if dir == "" {
		dir = "Logs"
	}
	exist, err := isPathExists(dir)
	if err != nil {
		fmt.Printf("get Logs/ dir error![%v]\n", err)
		return
	}
	if !exist {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir %s failed![%v]\n", dir, err)
		} else {
			fmt.Printf("mkdir %s success!\n", dir)
		}
	}

	if !boolStarted {
		for _, info := range mapLogInfo {
			go runPrint(info)
		}
	}
	boolStarted = true
}

func runPrint(info *logInfo) {
	lastDate := getNowTimeStr("2006_01_02")
	lastTime := time.Now().Unix() % 86400

	for true {
		if info.queue.Head() != nil { //日期变了要重新打开文件
			currDate := getNowTimeStr("2006_01_02")
			if lastDate != currDate {
				lastDate = currDate
			}
			fileName := fmt.Sprintf("%s/%s_%s_%d.zLog", dir, info.desc, lastDate, lastTime)
			if getFileSize(fileName) > maxFileSize {
				lastTime = time.Now().Unix() % 86400
				fileName = fmt.Sprintf("%s/%s_%s_%d.zLog", dir, info.desc, lastDate, lastTime)
			}
			logFile, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
			for info.queue.Head() != nil {
				msg := info.queue.Dequeue().(string)
				logger := log.New(logFile, "", log.LstdFlags)
				logger.Println(msg)
			}
			logFile.Close()
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// 输出字符串日志，带显示出控制台
func printLogger(str string, logLevel int) {
	if showFileLine {
		_, file, line, ok := runtime.Caller(2)
		if ok {
			str = "[" + file + " 行：" + strconv.Itoa(line) + "]：" + str
		}
	}
	logInfo := getLogInfo(logLevel)
	str = fmt.Sprintf("[%s] %s", logInfo.desc, str)
	if logLevel != LevelDebug {
		logInfo.queue.Enqueue(str)
	}
	if logLevel >= printLevel || logLevel == LevelDebug {
		color.Set(logInfo.color)
		fmt.Println(fmt.Sprintf("[%s]:%s", getNowTimeStr(""), str))
		color.Unset()
	}
}

func Debug(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelDebug)
}

func Log(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelLog)
}

func WarningLog(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelWarning)
}

func ErrorLog(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelError)
}

// 判断文件夹是否存在
func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func getFileSize(fileName string) int64 {
	fi, err := os.Stat(fileName)
	if err == nil {
		return fi.Size()
	}
	return 0
}

func getNowTimeStr(format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	timeObj := time.Now()
	return timeObj.Format(format)
}

func getLogInfo(logLevel int) *logInfo {
	return mapLogInfo[logLevel]
}
