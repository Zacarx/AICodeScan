package JDBC

import (
	"AICodeScan/FindFile"
	"fmt"
)

func FindJDBC(dir string) {
	FindFile.FindFileByJava(dir, "jdbc.txt", []string{"DriverManager.getConnection("})
	fmt.Println("JDBC分析完成")
}
