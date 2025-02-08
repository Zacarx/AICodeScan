package Unserialize

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Unserialize(dir string) {

	FindFile.FindFileByPHP(dir, "Unserialize.txt", []string{
		"__destruct(",
	})
	fmt.Println("PHP反序列化分析完成")

}
