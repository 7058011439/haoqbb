package Log

import (
	"fmt"
	"github.com/7058011439/haoqbb/File"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
	"github.com/gookit/color"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var printLevel int
var showFileLine bool
var maxFileSize int64
var dir string
var mutex sync.Mutex

type logLevel = int

const (
	LevelDebug logLevel = iota
	LevelLog
	LevelWarn
	LevelError
	LevelFatal
)

type logInfo struct {
	desc     string
	color    color.Color
	fileName string
}

func (l *logInfo) initFileName() {
	fileList, _ := File.WalkFile(dir, "", fmt.Sprintf("%v_%v", l.desc, getNowTimeStr("2006_01_02")))
	maxIndex := 0
	for _, fileName := range fileList {
		arr := strings.Split(fileName[:len(fileName)-5], "_")
		index := Util.StrToInt(arr[len(arr)-1])
		if index > maxIndex {
			maxIndex = index
		}
	}
	l.fileName = fmt.Sprintf("%v_%v_%v.zLog", l.desc, getNowTimeStr("2006_01_02"), maxIndex+1)
}

type logData struct {
	eType logLevel
	data  string
}

var mapLogInfo = map[int]*logInfo{
	LevelDebug: {desc: "Debug", color: color.Blue},
	LevelLog:   {desc: "Log", color: color.Green},
	LevelWarn:  {desc: "Warn", color: color.Yellow},
	LevelError: {desc: "Error", color: color.Red},
	LevelFatal: {desc: "Fatal", color: color.Red},
}

var queue = Stl.NewQueue()

func init() {
	Init(LevelLog, false, 0, "")
	for _, info := range mapLogInfo {
		info.initFileName()
	}
	go runPrint()
}

func Init(iPrintLevel int, bShowFileLine bool, iMaxFileSize int64, logDir string) {
	SetPrintLevel(iPrintLevel)
	SetShowFileLine(bShowFileLine)
	SetFileSize(iMaxFileSize)
	SetLogDir(logDir)
}

func SetPrintLevel(level int) {
	printLevel = level
}

func SetFileSize(iMaxFileSize int64) {
	if iMaxFileSize < 1 {
		maxFileSize = 5 * 1024 * 1024
	} else {
		maxFileSize = iMaxFileSize
	}
}

func SetLogDir(logDir string) {
	dir = logDir
	if dir == "" {
		dir = "Logs"
	}
	if !File.PathExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir %s failed![%v]\n", dir, err)
		} else {
			fmt.Printf("mkdir %s success!\n", dir)
		}
	}
}

func SetShowFileLine(bShowFileLine bool) {
	showFileLine = bShowFileLine
}

func print(level logLevel, data string) {
	info := mapLogInfo[level]
	info.color.Println(data)
}

func printFile() {
	mutex.Lock()
	defer mutex.Unlock()
	line := 1024 * 8
	// 这里连续两个for queue.Head() != nil。最外层的for是为了保证所有日志被输出; 里层的for是为了每次读取一定数量行(不全部读取),确保文件不会太大
	for queue.Head() != nil {
		typeData := map[int][]string{}
		typeData[LevelLog] = make([]string, 0, line)
		for queue.Head() != nil {
			msg := queue.Dequeue().(*logData)
			if len(typeData[msg.eType]) == 0 {
				typeData[msg.eType] = make([]string, 0, line)
			}
			typeData[msg.eType] = append(typeData[msg.eType], msg.data)

			// 任何类型日志都要出现在普通日志(Log)里面
			if msg.eType != LevelLog {
				typeData[LevelLog] = append(typeData[LevelLog], msg.data)
			}
			if len(typeData[LevelLog]) >= line {
				break
			}
		}
		for eType, dataList := range typeData {
			info := getLogInfo(eType)
			// 文件已达到指定大小, 或者过天了, 需要重新赋值文件名
			if File.GetFileSize(dir+"/"+info.fileName) > maxFileSize || strings.Index(info.fileName, getNowTimeStr("2006_01_02")) < 0 {
				info.initFileName()
			}
			logFile, _ := os.OpenFile(dir+"/"+info.fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
			for _, data := range dataList {
				logFile.WriteString(data + "\n")
			}
			logFile.Close()
		}
	}
}

func runPrint() {
	for {
		printFile()
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
	info := getLogInfo(logLevel)
	str = fmt.Sprintf("[%s] [%s] %s", getNowTimeStr(""), info.desc, str)
	// debug类型日志不存档
	if logLevel != LevelDebug {
		queue.Enqueue(&logData{
			eType: logLevel,
			data:  str,
		})
		if logLevel == LevelFatal {
			printFile()
		}
		if logLevel >= printLevel {
			print(logLevel, str)
		}
	} else {
		print(LevelDebug, str)
	}
}

func Debug(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelDebug)
}

func Log(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelLog)
}

func Warn(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelWarn)
}

func Error(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelError)
}

func Fatal(format string, args ...interface{}) {
	printLogger(fmt.Sprintf(format, args...), LevelFatal)
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
