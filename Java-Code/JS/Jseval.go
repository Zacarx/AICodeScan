package JS

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Eval(dir string) {
	fmt.Println("Eval分析开始")
	FindFile.FindFileByJava(dir, "eval.txt", []string{"eval("})
	fmt.Println("Eval分析完成")
}
