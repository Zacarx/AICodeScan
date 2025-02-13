package Reflect

import (
	"AICodeScan/FindFile"
	"fmt"
)

func ReflectCheck(dir string) {
	fmt.Println("反射分析开始")
	FindFile.FindFileByJava(dir, "fanshe.txt", []string{".invode("})
	fmt.Println("反射分析完成")
}
