package Fastjson

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Parsecheck(dir string) {
	fmt.Println("fastjson分析开始")
	FindFile.FindFileByJava(dir, "fastjson.txt", []string{".parseObject("})
	fmt.Println("fastjson分析完成")

}
