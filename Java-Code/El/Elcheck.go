package El

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Elcheck(dir string) {
	fmt.Println("表达式注入分析开始")
	//".getValue",  推荐不加
	FindFile.FindFileByJava(dir, "el.txt", []string{"SpelExpressionParser", "parseExpression"})
	fmt.Println("表达式注入分析完成")
}
