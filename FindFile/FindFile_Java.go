package FindFile

import (
	"AICodeScan/AItools"
	Rule2 "AICodeScan/CommonVul/Rule"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindFile FindFile函数用于在指定目录下查找符合规则的.java文件，并将包含规则的行写入到输出文件中
// 参数dir表示要搜索的目录路径
// 参数outputfile表示输出结果文件的路径（工具运行的目录）
// 参数rules表示要匹配的规则列表
func FindFileByJava(dir string, outputfile string, rules []string) {
	var fileList []string

	// 使用filepath.Walk遍历目标目录，跳过黑名单中的目录，收集所有.java文件的路径
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//如果f是一个文件夹
		if f.IsDir() {
			//继续进行遍历，如果在黑名单中的话就进行跳过
			if Rule2.MatchRule(path, Rule2.PathBlackJava) {
				return filepath.SkipDir
			}
			//如果文件存在的话就进行遍历 否则就进行判断，如果是java或者jsp后缀就添加到文件列表
		} else if strings.HasSuffix(f.Name(), ".java") || strings.HasSuffix(f.Name(), ".jsp") {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dir, err)
		return
	}

	// 检查遍历目录过程中的错误
	Check(err)

	// 创建或打开输出文件，以追加模式写入
	basedir := "./results/"

	err1 := os.MkdirAll(basedir, os.ModePerm)
	if err1 != nil {
		fmt.Println("Error creating directory:", err)
		return
	}
	outputfile = basedir + outputfile
	outputFile, err := os.OpenFile(outputfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	Check(err)
	defer outputFile.Close() // 确保文件在函数返回前被关闭
	for _, file := range fileList {
		f, err := os.Open(file)
		Check(err)

		// 使用defer关闭文件前，将其放入函数或代码块作用域中，以便及时释放资源
		func() {
			defer f.Close()

			scanner := bufio.NewScanner(f)
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 10*1024*1024)

			lineNumber := 1
			var lastFile string
			fileProcessed := false // 标记当前文件是否已被处理

			for scanner.Scan() {
				if fileProcessed {
					break // 如果文件已经调用过AI分析，跳出循环
				}

				line := strings.TrimSpace(scanner.Text())
				for _, rule := range rules {
					if strings.Contains(strings.ToLower(line), strings.ToLower(rule)) {
						if Rule2.MatchRule(line, Rule2.LineBlack) {
							break
						}
						if !Rule2.RemoveStaticVar(strings.ToLower(line), strings.ToLower(rule)) {
							break
						}

						if lastFile != file {
							_, err := outputFile.WriteString(fmt.Sprintf("====================================================================\n\n"))
							result, err := AItools.AWA(file, lineNumber, line)
							if err != nil {
								fmt.Println("Error in AI analysis:", err)
								continue
							}
							fmt.Print(result)
							//_, err = outputFile.WriteString(fmt.Sprintf("AI Analysis Result: %s\n\n", result))
							_, err = outputFile.WriteString(fmt.Sprintf("file [%s]\n%d : %s\n\n AI Analysis Result:\n%s\n\n", file, lineNumber, line, result))
							if err != nil {
								fmt.Println("Error writing to file:", err)
							}
						}
						lastFile = file
					}
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
		}()
	}
}
