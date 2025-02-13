package SSRF

import (
	"AICodeScan/FindFile"
	"fmt"
)

func PHP_SSRF(dir string) {
	fmt.Println("SSRF分析开始")
	FindFile.FindFileByPHP(dir, "SSRF.txt", []string{
		"curl_exec(",
	})
	fmt.Println("SSRF分析完成")
}
