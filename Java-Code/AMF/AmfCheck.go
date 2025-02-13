package AMF

import (
	"AICodeScan/FindFile"
	"fmt"
)

func AmfCheck(dir string) {
	fmt.Println("AMF检查开始")
	FindFile.FindFileByJava(dir, "AmfCheck.txt", []string{".readMessage("})
	fmt.Println("AMF检查完成")

}
