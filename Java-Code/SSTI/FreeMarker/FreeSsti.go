package FreeMarker

import (
	"AICodeScan/FindFile"
	"fmt"
)

func FreeSsti(dir string) {
	FindFile.FindFileByJava(dir, "Freemarkssti.txt", []string{"new Template("})
	fmt.Println("FreeMarker SSTI 分析完成")
}
