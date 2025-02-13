package FileWrite

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Write(dir string) {
	fmt.Println("PHP文件写入分析开始")
	FindFile.FindFileByPHP(dir, "FileWrite.txt", []string{
		"file_put_contents(",
	})
	fmt.Println("PHP文件写入分析完成")

}
