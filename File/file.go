package File

import (
	"bufio"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func ReadLine(filePath string, hookFn func(string, int, string, ...interface{}), args ...interface{}) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	line := 1
	bfRd := bufio.NewReader(f)
	codeType := ""
	fileData, _ := ioutil.ReadFile(filePath)
	detector := chardet.NewTextDetector()
	if charset, err := detector.DetectBest(fileData); err == nil {
		codeType = charset.Charset
	}
	for {
		msg, _, err := bfRd.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		data := string(msg)

		if codeType == "UTF-8" {
			hookFn(filePath, line, data, args...)
		} else {
			hookFn(filePath, line, mahonia.NewDecoder("GBK").ConvertString(data), args...)
		}
		line++
	}
	return nil
}

func Split(r rune) bool {
	return r == ' ' || r == '|' || r == '\t' || r == '/'
}

func WalkFile(dirPth, suffix string, prefix string) (files []string, err error) {
	_, err = os.Stat(dirPth)
	if err != nil {
		fmt.Printf("遍历目录失败, 不存在目录:%v\n", dirPth)
		return
	}

	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	prefix = strings.ToUpper(prefix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) && strings.HasPrefix(strings.ToUpper(fi.Name()), prefix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func WalkPath(dirPth, suffix string, prefix string) (files []string, err error) {
	_, err = os.Stat(dirPth)
	if err != nil {
		fmt.Printf("遍历目录失败, 不存在目录:%v\n", dirPth)
		return
	}

	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	prefix = strings.ToUpper(prefix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() {
			if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) && strings.HasPrefix(strings.ToUpper(fi.Name()), prefix) {
				files = append(files, filename)
			}
			return nil
		}
		return nil
	})
	return files, err
}

func CreateFileWithDir(path string, name string) *os.File {
	fullPath := path + "\\" + name
	if path == "" || name == "" {
		fullPath = path + name
	}
	fullPath = strings.ReplaceAll(fullPath, "\\\\", "\\")
	pos := strings.LastIndex(fullPath, "\\")
	os.MkdirAll(fullPath[0:pos], os.ModePerm)
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
	}
	return file
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	pos := strings.LastIndexByte(dst, '\\')
	poso := strings.LastIndexByte(dst, '/')
	fmt.Println(src, dst, poso, pos)

	if poso > pos {
		pos = poso
	}
	dstPath := dst[:pos]
	dstName := dst[pos+1:]

	destination := CreateFileWithDir(dstPath, dstName)
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyFile(src, des string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		fmt.Println(err)
	}
	defer desFile.Close()

	return io.Copy(desFile, srcFile)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	//当为空文件或文件夹存在
	if err == nil {
		return true
	}
	//os.IsNotExist(err)为true，文件或文件夹不存在
	if os.IsNotExist(err) {
		return false
	}
	//其它类型，不确定是否存在
	return false
}

func GetFileSize(fileName string) int64 {
	fi, err := os.Stat(fileName)
	if err == nil {
		return fi.Size()
	}
	return 0
}

// extName: 是否包含扩展名
func GetFileName(path string, extName bool) string {
	fileNameWithExt := filepath.Base(path)
	if extName {
		return fileNameWithExt
	} else {
		return strings.TrimSuffix(fileNameWithExt, filepath.Ext(path))
	}
}

// CreateDir 根据文件名(也有可能是目录)创建目录，如果目录存在，不创建
// 如
// fileName             类型     文件/目录是否存在         处理结果
// F:\Haoqbb\Logs     	目录     是                      不处理
// F:\Haoqbb\Logs     	目录     否                      创建目录:F:\Haoqbb\Logs
// F:\Haoqbb\Logs\1.log	文件     是                      不处理
// F:\Haoqbb\Logs\2.log	文件     否                      创建目录:F:\Haoqbb\Logs
func CreateDir(fileName string) error {
	if !PathExists(fileName) {
		// 获取目录部分
		dir := fileName
		if fileInfo, err := os.Stat(fileName); err == nil && !fileInfo.IsDir() {
			// 如果路径存在且是文件，则使用文件的目录部分
			dir = filepath.Dir(fileName)
		} else if filepath.Ext(fileName) != "" {
			// 如果路径有扩展名，但不存在，则假设是文件，获取目录部分
			dir = filepath.Dir(fileName)
		}
		return os.MkdirAll(dir, os.ModePerm) // 创建目录
	}
	return nil
}
