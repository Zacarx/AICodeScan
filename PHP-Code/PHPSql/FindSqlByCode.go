package PHPSql

import (
	"AICodeScan/AItools"
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// 函数用于检查是否存在php代码内容，并将相关信息写入 sql.txt
func findSqlByCode(dir string) {
	selectList := []string{}
	var lastFile string // 记录上一次输出的文件，用于控制输出格式

	keywords := []string{"like '%\" +", "mysql_query(", "->where(", "->order(", "mysqli_query("}

	err := filepath.Walk(dir, func(path string, f fs.FileInfo, err error) error {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".php") {
			// 打开文件
			lineNumber := 1 // 行号，用于标识匹配行的位置
			file, err := os.Open(path)
			check(err)
			defer file.Close()

			// 逐行扫描文件内容
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				// 如果行中包含 @Select 注解，则将相关信息添加到 selectList 中
				for _, keyword := range keywords {
					if strings.Contains(line, keyword) {
						if lastFile != file.Name() {
							selectList = append(selectList, fmt.Sprintf("====================================================================\n"))
							//selectList = append(selectList, fmt.Sprintf("file [%s]\n%d: %s", file.Name(), lineNumber, line))
							filename := file.Name()
							result, err := AItools.AWA(filename, lineNumber, line)
							//fmt.Println(result)
							if err != nil {
								fmt.Println("Error in AI analysis:", err)
							}
							selectList = append(selectList, fmt.Sprintf("file [%s]\n%d: %s\nAI1 Analysis Result: %s\n\n", filename, lineNumber, line, result))
							lastFile = file.Name()
						} else {
							selectList = append(selectList, fmt.Sprintf("====================================================================\n"))
							selectList = append(selectList, fmt.Sprintf("%d : %s", lineNumber, line))
							result, err := AItools.AWA(file.Name(), lineNumber, line)
							if err != nil {
								fmt.Println("Error in AI analysis:", err)
								continue
							}
							selectList = append(selectList, fmt.Sprintf("AI Analysis Result: %s\n\n", result))

						}
					}

				}
				lineNumber++
			}

		}
		return nil
	})
	check(err)

	// 如果存在 @Select 注解，则将相关信息写入到 sql.txt 文件中
	if len(selectList) > 0 {
		writeToFile("sql.txt", selectList)
	}
}
