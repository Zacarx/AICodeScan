package main

import (
	"AICodeScan/Utils"
	"fmt"
	"github.com/common-nighthawk/go-figure"
)

func main() {
	myFigure := figure.NewFigure("AICodeScan", "", true)
	myFigure.Print()
	fmt.Println("									by Zacarx")

	Utils.Start()
	//elapsed := time.Since()      // 计
	//color.Green("[+] 扫描完成! 花费时长:%s\n", elapsed) // 算经过的时间
}
