package JNDI

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Jndi(dir string) {
	fmt.Println("JNDI分析开始")
	FindFile.FindFileByJava(dir, "jndi.txt", []string{".lookup("})
	fmt.Println("JNDI分析完成")
}
