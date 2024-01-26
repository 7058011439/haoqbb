package Log

import (
	"fmt"
	"github.com/7058011439/haoqbb/File"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
	"github.com/fatih/color"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var printLevel int
var showFileLine bool
var maxFileSize int64
var dir string

type logLevel = int

const (
	LevelDebug   logLevel = 0
	LevelLog     logLevel = 1
	LevelWarning logLevel = 2
	LevelError   logLevel = 3
)

type logInfo struct {
	desc     string
	color    *color.Color
	fileName string
}

type logData struct {
	eType logLevel
	data  string
}

var mapLogInfo = map[int]*logInfo{
	LevelDebug:   {desc: "Debug", color: color.New(color.FgBlue)},
	LevelLog:     {desc: "Log", color: color.New(color.FgGreen)},
	LevelWarning: {desc: "Warning", color: color.New(color.FgYellow)},
	LevelError:   {desc: "Error", color: color.New(color.FgRed)},
}

var queue = Stl.NewQueue()

func init() {
	Init(LevelLog, false, 0, "")
	go runPrint()
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
	if !File.PathExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir %s failed![%v]\n", dir, err)
		} else {
			fmt.Printf("mkdir %s success!\n", dir)
		}
	}
	for eType, info := range mapLogInfo {
		info.fileName = newFileName(eType)
	}
}

func newFileName(eType int) string {
	info := getLogInfo(eType)
	fileList, _ := File.WalkFile(dir, "", fmt.Sprintf("%v_%v", info.desc, getNowTimeStr("2006_01_02")))
	maxIndex := 0
	for _, fileName := range fileList {
		arr := strings.Split(fileName[:len(fileName)-5], "_")
		index := Util.StrToInt(arr[len(arr)-1])
		if index > maxIndex {
			maxIndex = index
		}
	}
	return fmt.Sprintf("%v_%v_%v.zLog", info.desc, getNowTimeStr("2006_01_02"), maxIndex+1)
}

func runPrint() {
	for true {
		if queue.Head() != nil {
			typeData := map[int][]string{}
			for queue.Head() != nil {
				msg := queue.Dequeue().(*logData)
				typeData[msg.eType] = append(typeData[msg.eType], msg.data)

				// 任何类型日志都要出现在普通日志(Log)里面
				if msg.eType != LevelLog {
					typeData[LevelLog] = append(typeData[LevelLog], msg.data)
				}
				if msg.eType >= printLevel || msg.eType == LevelDebug {
					info := mapLogInfo[msg.eType]
					info.color.EnableColor()
					info.color.Set()
					fmt.Println(msg.data)
					info.color.DisableColor()
				}
			}
			for eType, dataList := range typeData {
				info := getLogInfo(eType)
				if File.GetFileSize(dir+"/"+info.fileName) > maxFileSize || strings.Index(info.fileName, getNowTimeStr("2006_01_02")) < 0 {
					info.fileName = newFileName(eType)
				}
				logFile, _ := os.OpenFile(dir+"/"+info.fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
				for _, data := range dataList {
					logFile.WriteString(data + "\n")
				}
				logFile.Close()
			}
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
	info := getLogInfo(logLevel)
	str = fmt.Sprintf("[%s] [%s] %s", getNowTimeStr(""), info.desc, str)
	if logLevel != LevelDebug {
		queue.Enqueue(&logData{
			eType: logLevel,
			data:  str,
		})
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
