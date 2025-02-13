package FindFile

import (
	"AICodeScan/AItools"
	Rule2 "AICodeScan/CommonVul/Rule"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

// FindFileByJava 函数用于在指定目录下查找符合规则的 .java 文件，并将包含规则的行写入到输出文件中
// 参数 dir 表示要搜索的目录路径
// 参数 outputfile 表示输出结果文件的路径（工具运行的目录）
// 参数 rules 表示要匹配的规则列表
func FindFileByJava(dir string, outputfile string, rules []string) {
	var fileList []string

	// 统计符合条件的文件数量
	fileCount := 0
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			if Rule2.MatchRule(path, Rule2.PathBlackJava) {
				return filepath.SkipDir
			}
		} else if strings.HasSuffix(f.Name(), ".java") || strings.HasSuffix(f.Name(), ".jsp") {
			fileCount++
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dir, err)
		return
	}

	if fileCount == 0 {
		fmt.Println("No .java or .jsp files found.")
		return
	}

	// 创建进度条
	bar := pb.StartNew(fileCount)

	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			if Rule2.MatchRule(path, Rule2.PathBlackJava) {
				return filepath.SkipDir
			}
		} else if strings.HasSuffix(f.Name(), ".java") || strings.HasSuffix(f.Name(), ".jsp") {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dir, err)
		return
	}

	basedir := "./results/"
	err1 := os.MkdirAll(basedir, os.ModePerm)
	if err1 != nil {
		fmt.Println("Error creating directory:", err1)
		return
	}
	outputfile = basedir + outputfile
	outputFile, err := os.OpenFile(outputfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening output file:", err)
		return
	}
	defer outputFile.Close()

	for _, file := range fileList {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("Error opening file:", err)
			bar.Increment()
			continue
		}

		func() {
			defer f.Close()

			scanner := bufio.NewScanner(f)
			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, 10*1024*1024)

			lineNumber := 1
			var lastFile string
			fileProcessed := false

			for scanner.Scan() {
				if fileProcessed {
					break
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
							_, err = outputFile.WriteString(fmt.Sprintf("file [%s]\n%d : %s\n\n AI Analysis Result:\n%s\n\n", file, lineNumber, line, result))
							if err != nil {
								fmt.Println("Error writing to file:", err)
							}
						}
						lastFile = file
						fileProcessed = true
					}
				}
				lineNumber++
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "reading standard input:", err)
			}
		}()

		bar.Increment()
	}

	bar.Finish()
}
