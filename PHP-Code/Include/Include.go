package Include

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Include(dir string) {
	fmt.Println("PHP文件包含分析开始")
	FindFile.FindFileByPHP(dir, "Include.txt", []string{
		"include(",
	})
	fmt.Println("PHP文件包含分析完成")
}
