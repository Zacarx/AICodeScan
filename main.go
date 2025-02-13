package main

import (
	"AICodeScan/Utils"
	"AICodeScan/output"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"strings"
	"time"
)

func main() {
	// 打印炫酷标题
	myFigure := figure.NewFigure("AICodeScan", "", true)
	printGradientFigure(myFigure.String())
	//fmt.Println()
	printGradient("\t\t\t\t\t\t\t\tby Zacarx")

	start := time.Now()
	// 启动扫描
	Utils.Start()
	output.OT()

	// 计算扫描时间
	defer func() {
		elapsed := time.Since(start)
		printGradient("扫描完成! 花费时长: %s\n", elapsed)
	}()
}

// printGradientFigure 打印具有渐变颜色的文本（适用于 figure 字符串）
func printGradientFigure(text string) {
	// 定义渐变颜色序列
	colors := []color.Attribute{
		color.FgHiGreen,
		color.FgGreen,
		color.FgYellow,
		color.FgHiYellow,
		color.FgHiRed,
		color.FgRed,
	}

	// 获取每行文本
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		runes := []rune(line)
		//runesLength := utf8.RuneCountInString(line)

		for i, r := range runes {
			idx := i % len(colors)
			c := color.New(colors[idx])
			c.Print(string(r))
		}
		fmt.Println()
	}
}

// printGradient 打印具有渐变颜色的文本
func printGradient(format string, a ...interface{}) {
	// 定义渐变颜色序列
	colors := []color.Attribute{
		color.FgHiGreen,
		color.FgGreen,
		color.FgYellow,
		color.FgHiYellow,
		color.FgHiRed,
		color.FgRed,
	}

	// 获取格式化后的文本
	text := fmt.Sprintf(format, a...)
	runes := []rune(text)
	//runesLength := utf8.RuneCountInString(text)

	for i, r := range runes {
		idx := i % len(colors)
		c := color.New(colors[idx])
		c.Print(string(r))
	}
	fmt.Println()
}

// 如果需要使用中文或特殊字符，确保正确处理
//func utf8.RuneCountInString(s string) int {
//	return utf8.RuneCountInString(s)
//}
