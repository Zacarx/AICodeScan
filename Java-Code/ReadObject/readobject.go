package ReadObject

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Readobjectcheck(dir string) {
	fmt.Println("反序列化分析开始")
	FindFile.FindFileByJava(dir, "readobject.txt", []string{".readobject(", ".deserialize("})
	fmt.Println("反序列化分析完成")
}
