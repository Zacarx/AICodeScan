package FreeMarker

import (
	"AICodeScan/FindFile"
	"fmt"
)

func FreeSsti(dir string) {
	fmt.Println("FreeMarker SSTI 分析开始")
	FindFile.FindFileByJava(dir, "Freemarkssti.txt", []string{"new Template("})
	fmt.Println("FreeMarker SSTI 分析完成")
}
