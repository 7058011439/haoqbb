package File

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ZipFile(dest string, src string) error {
	zfile, err := os.Create(dest)
	if err != nil {
		fmt.Printf("创建压缩文件失败，err = %v\n", err)
		return err
	}
	defer zfile.Close()
	zipWriter := zip.NewWriter(zfile)
	defer zipWriter.Close()

	f1, err := os.Open(src)
	if err != nil {
		fmt.Printf("压缩文件失败, 打开文件错误, err = %v\n", err)
		return err
	}
	defer f1.Close()

	shortName := filepath.Base(src)
	w1, err := zipWriter.Create(shortName)
	if err != nil {
		fmt.Printf("创建压缩信息失败, err = %v\n", err)
		return err
	}
	if _, err := io.Copy(w1, f1); err != nil {
		fmt.Printf("压缩内容失败, err = %v\n", err)
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
