package ReadObject

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Readobjectcheck(dir string) {
	FindFile.FindFileByJava(dir, "readobject.txt", []string{".readobject(", ".deserialize("})
	fmt.Println("反序列化分析完成")
}
