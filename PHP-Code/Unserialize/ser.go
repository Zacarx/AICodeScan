package Unserialize

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Unserialize(dir string) {
	fmt.Println("PHP反序列化分析开始")
	FindFile.FindFileByPHP(dir, "Unserialize.txt", []string{
		"__destruct(",
	})
	fmt.Println("PHP反序列化分析完成")

}
