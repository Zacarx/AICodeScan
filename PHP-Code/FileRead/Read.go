package FileRead

import (
	"AICodeScan/CommonVul/Rule"
	"AICodeScan/FindFile"
	"fmt"
)

func Read(dir string) {
	fmt.Println("PHP文件读取分析开始")
	FindFile.FindFileByPHP(dir, "FileRead_Phar.txt", Rule.PHPFileReadList)
	fmt.Println("PHP文件读取分析完成")

}
