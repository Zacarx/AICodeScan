package Rce

import (
	"AICodeScan/CommonVul/Rule"
	"AICodeScan/FindFile"
	"fmt"
)

func JavaRce(dir string) {
	fmt.Println("RCE分析开始")
	FindFile.FindFileByJava(dir, "rce.txt", Rule.JavaRceRuleList)
	fmt.Println("RCE分析完成")
}

func PHPRce(dir string) {
	fmt.Println("RCE分析开始")
	FindFile.FindFileByPHP(dir, "rce.txt", Rule.PHPRceRuleList)
	fmt.Println("RCE分析完成")
}
