package JS

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Eval(dir string) {
	FindFile.FindFileByJava(dir, "eval.txt", []string{"eval("})
	fmt.Println("Eval分析完成")
}
