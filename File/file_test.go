package File

import "testing"

func TestZip(t *testing.T) {
	//ZipDir("D:\\桌面\\2022_11_01.zip", "D:\\桌面\\2022_11_01")
	//ZipDirGZ("D:\\桌面\\2022_11_01.rar.gz", "D:\\桌面\\2022_11_01")
	ZipDir7z("D:\\桌面\\2022_11_01.7z", "D:\\桌面\\2022_11_01")
}

func TestPathExists(t *testing.T) {
	testCases := []struct {
		fileName string
		expected bool
	}{
		{"F:\\Haoqbb\\Go\\Core\\File", true},           // 存在目录(正斜线)
		{"F:/Haoqbb/Go/Core/File", true},               // 存在目录(反斜线)
		{"F:\\Haoqbb\\Go\\Core\\Fil", false},           // 不存在目录(正斜线)
		{"F:/Haoqbb/Go/Core/Fil", false},               // 不存在目录(反斜线)
		{"F:\\Haoqbb\\Go\\Core\\File\\7z.dll", true},   // 存在文件(正斜线,带后缀名)
		{"F:/Haoqbb/Go/Core/File/7z.dll", true},        // 存在文件(反斜线,带后缀名)
		{"F:\\Haoqbb\\Go\\Core\\File\\7z.dl", false},   // 不存在文件(正斜线,带后缀名)
		{"F:/Haoqbb/Go/Core/File/7z.dl", false},        // 不存在文件(反斜线,带后缀名)
		{"F:\\Haoqbb\\Go\\Core\\File\\testdata", true}, // 存在文件(正斜线,不带后缀名)
		{"F:/Haoqbb/Go/Core/File/testdata", true},      // 存在文件(正斜线,不带后缀名)
		{"F:\\Haoqbb\\Go\\Core\\File\\test", false},    // 不存在文件(正斜线,不带后缀名)
		{"F:/Haoqbb/Go/Core/File/test", false},         // 不存在文件(反斜线,不带后缀名)
	}

	for _, testCase := range testCases {
		actualRet := PathExists(testCase.fileName)
		if actualRet != testCase.expected {
			t.Errorf("测试失败, 期望结果 = %v, 实际结果 = %v, 测试数据 = %v", testCase.expected, actualRet, testCase.fileName)
		}
	}
}
