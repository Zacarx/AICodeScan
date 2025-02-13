package Auth_Bypass

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Auth(dir string) {
	fmt.Println("权限绕过分析开始")
	FindFile.FindFileByJava(dir, "Auth_Bypass.txt", []string{".getRequestURL(", ".getRequestURI("})
	fmt.Println("权限绕过分析完成")

}
