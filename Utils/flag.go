package Utils

import (
	Rule2 "AICodeScan/CommonVul/Rule"
	"AICodeScan/Filter"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"strings"
)

var (
	Dir      *string
	language *string
	help     *string
)

func Start() {
	// 开始审计

	parseFlag()
	*language = strings.ToLower(*language)
	if *language == "java" {
		Java_Codeing()
	}

	if *language == "php" {
		PHP_Codeing()
	}

}

func parseFlag() {

	// 高级命令行解析
	help = flag.String("h", "", "使用帮助")
	Dir = flag.String("d", "", "要扫描的目录")
	language = flag.String("L", "", "审计语言")
	pathBlackRule := flag.String("pb", "", "路径黑名单")
	lineBlackRule := flag.String("lb", "", "行黑名单")
	uploadRule := flag.String("u", "", "文件上传规则")
	rceRule := flag.String("r", "", "RCE规则")
	filterfile := flag.String("m", "", "过滤的字符串")
	//flag.String("o", "./results/results.html", "Output HTML file path")
	//outdir := flag.String("o", "", "输出结果")
	flag.Parse()

	if *language == "" && *filterfile == "" {
		color.Red("请使用 -H 选项提供帮助")
		return
	}

	if *language != "" {

		if *Dir != "" {
			*Dir = ClearDir(*Dir)
			if *pathBlackRule != "" {
				// 读取路径黑名单
				Rule2.PathBlackJava = append(Rule2.PathBlackJava, *pathBlackRule)
				fmt.Println("路径黑名单:", Rule2.PathBlackJava)
			} // 所有要执行的扫描函数

			if *lineBlackRule != "" {
				Rule2.LineBlack = append(Rule2.LineBlack, *lineBlackRule)
			}

			if *uploadRule != "" {
				if *language == "java" {
					Rule2.JavaUploadRuleList = append(Rule2.JavaUploadRuleList, *uploadRule)
				} else if *language == "php" {
					Rule2.PHPUploadRuleList = append(Rule2.PHPUploadRuleList, *uploadRule)
				}

			}

			if *rceRule != "" {
				Rule2.JavaRceRuleList = append(Rule2.JavaRceRuleList, *rceRule)
			}

		}

	}

	if *filterfile != "" {
		if *Dir != "" {
			Filter.FilterFile(*filterfile, *Dir)

		} else {
			color.Red("请使用 -d 选项提供目录")
			return
		}

	}

}
