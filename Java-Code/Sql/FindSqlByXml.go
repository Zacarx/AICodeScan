package Sql

import (
	"AICodeScan/AItools"
	Rule2 "AICodeScan/CommonVul/Rule"
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// findSqlByXml 函数用于检查 XML 文件中的关键字并进行 AI 分析
func findSqlByXml(dir string) {
	xmlList := []string{}
	var lastFile string // 记录上一次输出的文件，用于控制输出格式

	// 使用 Walk 函数遍历目录，查找所有的 .xml 文件
	err := filepath.Walk(dir, func(path string, f fs.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".xml") {
			// xml黑名单匹配
			if Rule2.MatchRule(f.Name(), Rule2.XmlBlack) {
				return nil
			}
			xmlList = append(xmlList, path)
		}
		return nil
	})
	check(err)

	// 定义需要搜索的关键字
	keywords := []string{"${", "like '%${", "order by ${"} // 这里可以添加更多关键字

	// 遍历 XML 文件列表
	for _, file := range xmlList {
		f, err := os.Open(file)
		check(err)

		func() {
			defer f.Close()

			scanner := bufio.NewScanner(f)
			lineNumber := 1
			fileProcessed := false // 标记当前文件是否已被处理
			foundKeywords := []string{}

			for scanner.Scan() {
				if fileProcessed {
					break // 如果文件已经调用过AI分析，跳出循环
				}

				line := strings.TrimSpace(scanner.Text())
				if Rule2.MatchRule(line, Rule2.XmlSqlBlack) {
					continue
				}

				for _, keyword := range keywords {
					if strings.Contains(line, keyword) {
						if lastFile != file {
							foundKeywords = append(foundKeywords, fmt.Sprintf("====================================================================\n"))
							result, err := AItools.AWA(file, lineNumber, line)
							if err != nil {
								fmt.Println("Error in AI analysis:", err)
								continue
							}
							foundKeywords = append(foundKeywords, fmt.Sprintf("file [%s]\n%d: %s\nAI Analysis Result:\n%s\n\n", file, lineNumber, line, result))
							lastFile = file
						} else {
							foundKeywords = append(foundKeywords, fmt.Sprintf("====================================================================\n"))
							result, err := AItools.AWA(file, lineNumber, line)
							if err != nil {
								fmt.Println("Error in AI analysis:", err)
								continue
							}
							foundKeywords = append(foundKeywords, fmt.Sprintf("%d : %sAI Analysis Result:\n%s\n\n", lineNumber, line, result))
						}

						// 调用 AI 审计函数

					}
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}

			// 如果找到关键字，则将相关信息写入到 sql.txt 文件中
			if len(foundKeywords) > 0 {
				writeToFile("sql.txt", foundKeywords)
			}
		}()
	}
}
