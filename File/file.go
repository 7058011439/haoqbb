package File

import (
	"Core/Log"
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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
		Log.ErrorLog("遍历目录失败, 不存在目录:%v", dirPth)
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
		Log.ErrorLog("遍历目录失败, 不存在目录:%v", dirPth)
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

func ZipFile(dest string, src string) error {
	zfile, err := os.Create(dest)
	if err != nil {
		Log.ErrorLog("创建压缩文件失败，err = %v", err)
		return err
	}
	defer zfile.Close()
	zipWriter := zip.NewWriter(zfile)
	defer zipWriter.Close()

	f1, err := os.Open(src)
	if err != nil {
		Log.ErrorLog("压缩文件失败, 打开文件错误, err = %v", err)
		return err
	}
	defer f1.Close()

	shortName := filepath.Base(src)
	w1, err := zipWriter.Create(shortName)
	if err != nil {
		Log.ErrorLog("创建压缩信息失败, err = %v", err)
		return err
	}
	if _, err := io.Copy(w1, f1); err != nil {
		Log.ErrorLog("压缩内容失败, err = %v", err)
		return err
	}
	return nil
}

func ZipDir7z(zipPath string, path string) error {
	cmd := exec.Command("7z.exe", "a", "-mx4", zipPath, path)
	//cmd.Stdout = os.Stdout
	return cmd.Run()
}

func ZipDirGZ(zipPath string, path string) error {
	// create zip file
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// new zip writer
	zipWriter, _ := gzip.NewWriterLevel(archive, gzip.BestCompression)
	defer zipWriter.Close()

	tw := tar.NewWriter(zipWriter)
	defer tw.Close()

	// visit all the files or directories in the tree
	err = filepath.Walk(path, func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		head, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		head.Name, err = filepath.Rel(filepath.Dir(path), fileName)

		if err = tw.WriteHeader(head); err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		io.Copy(tw, f)
		return err
	})
	return err
}

func ZipDir(zipPath string, path string) error {
	// create zip file
	if err := os.MkdirAll(filepath.Dir(zipPath), os.ModePerm); err != nil {
		return err
	}
	archive, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// new zip writer
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	// visit all the files or directories in the tree
	err = filepath.Walk(path, func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(path), fileName)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		}

		// create writer for the file header and save content of the file
		headerWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(fileName)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(headerWriter, f)
		return err
	})
	return err
}

func UnZip(zipFileName string, destDir string) error {
	archive, err := zip.OpenReader(zipFileName)
	if err != nil {
		return fmt.Errorf("读取压缩文件错误, err = %v", err)
	}

	defer archive.Close()

	if destDir == "" {
		destDir, _ = filepath.Split(zipFileName)
	}

	for _, f := range archive.File {
		filePath := filepath.Join(destDir, f.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path = %v", destDir)
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败, err = %v", err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("创建文件失败, err = %v", err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("打开文件失败, err = %v", err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("复制类容失败, err = %v", err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
