package File

import "testing"

func TestZip(t *testing.T) {
	//ZipDir("D:\\桌面\\2022_11_01.zip", "D:\\桌面\\2022_11_01")
	//ZipDirGZ("D:\\桌面\\2022_11_01.rar.gz", "D:\\桌面\\2022_11_01")
	ZipDir7z("D:\\桌面\\2022_11_01.7z", "D:\\桌面\\2022_11_01")
}
