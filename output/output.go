package output

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Vulnerability struct {
	FilePath    string
	Line        string
	Code        string
	VulnType    string
	Severity    string
	Description string
	Payload     string
}

type FileVulnerabilities struct {
	FileName        string
	Vulnerabilities []Vulnerability
}

func OT() {
	// 命令行参数解析，输入结果保存路径
	outputPath := "./results/results.html"

	files, err := filepath.Glob("./results/*.txt")
	if err != nil {
		fmt.Println("Failed to read result files:", err)
		return
	}

	fileVulnerabilitiesMap := make(map[string]*FileVulnerabilities)
	for _, file := range files {
		contents, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println("Failed to read file:", file)
			continue
		}

		fileName := filepath.Base(file)
		fileVulnerabilities, ok := fileVulnerabilitiesMap[fileName]
		if !ok {
			fileVulnerabilities = &FileVulnerabilities{FileName: fileName}
			fileVulnerabilitiesMap[fileName] = fileVulnerabilities
		}

		fileVulnerabilities.Vulnerabilities = append(fileVulnerabilities.Vulnerabilities, parseVulnerabilities(contents)...)
	}

	fileVulnerabilitiesList := make([]FileVulnerabilities, 0, len(fileVulnerabilitiesMap))
	for _, fv := range fileVulnerabilitiesMap {
		fileVulnerabilitiesList = append(fileVulnerabilitiesList, *fv)
	}

	if err := generateHTML(fileVulnerabilitiesList, outputPath); err != nil {
		fmt.Println("Failed to generate HTML output:", err)
	}
}

func parseVulnerabilities(contents []byte) []Vulnerability {
	var vulnerabilities []Vulnerability
	lines := strings.Split(string(contents), "\n")
	var currentVuln *Vulnerability
	inAnalysis := false
	inPayload := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "====================================================================" {
			if currentVuln != nil && (currentVuln.Code != "" || currentVuln.VulnType != "") {
				if !strings.Contains(currentVuln.Description, "大概率不存在漏洞") && !strings.Contains(currentVuln.Description, "不存在漏洞") {
					vulnerabilities = append(vulnerabilities, *currentVuln)
				}
			}
			currentVuln = nil
			inAnalysis = false
			inPayload = false
		} else if strings.HasPrefix(line, "file [") && strings.HasSuffix(line, "]") {
			if currentVuln != nil {
				if !strings.Contains(currentVuln.Description, "大概率不存在漏洞") && !strings.Contains(currentVuln.Description, "不存在漏洞") {
					vulnerabilities = append(vulnerabilities, *currentVuln)
				}
			}
			filePath := strings.TrimSuffix(strings.TrimPrefix(line, "file ["), "]")
			currentVuln = &Vulnerability{FilePath: filePath}
		} else if currentVuln != nil {
			if strings.HasPrefix(line, "AI Analysis Result:") {
				inAnalysis = true
			} else if !inAnalysis && strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					lineNum := strings.TrimSpace(parts[0])
					code := strings.TrimSpace(parts[1])
					currentVuln.Line = lineNum
					currentVuln.Code = code
				}
			} else if inAnalysis {
				if strings.HasPrefix(line, "漏洞类型：") {
					parts := strings.SplitN(line, "：", 2)
					if len(parts) == 2 {
						currentVuln.VulnType = strings.TrimSpace(parts[1])
					}
				} else if strings.HasPrefix(line, "危害等级：") {
					parts := strings.SplitN(line, "：", 2)
					if len(parts) == 2 {
						currentVuln.Severity = strings.TrimSpace(parts[1])
					}
				} else if strings.HasPrefix(line, "判断理由：") {
					parts := strings.SplitN(line, "：", 2)
					if len(parts) == 2 {
						currentVuln.Description = strings.TrimSpace(parts[1])
					}
				} else if strings.HasPrefix(line, "payload：") {
					currentVuln.Payload = strings.TrimSpace(strings.TrimPrefix(line, "payload："))
					inPayload = true
				} else if inPayload && len(line) > 0 {
					currentVuln.Payload = strings.TrimSpace(currentVuln.Payload + "\n" + line)
				}
			}
		}
	}

	// 添加最后一个漏洞
	if currentVuln != nil && (currentVuln.Code != "" || currentVuln.VulnType != "") {
		if !strings.Contains(currentVuln.Description, "大概率不存在漏洞") && !strings.Contains(currentVuln.Description, "不存在漏洞") {
			vulnerabilities = append(vulnerabilities, *currentVuln)
		}
	}
	return vulnerabilities
}

func generateHTML(data []FileVulnerabilities, path string) error {
	t := template.Must(template.New("report").Parse(`
<!DOCTYPE html>
<html>
<head>
<title>Vulnerability Report</title>
<meta charset="UTF-8">
<style>
body {
    font-family: 'Fira Code', monospace;
    background-color: #282c34;
    color: #abb2bf;
    padding: 20px;
    transition: background-color 0.3s ease;
}
h1 {
    color: #61afef;
    text-align: center;
    margin-bottom: 30px;
    font-size: 3em;
    transition: color 0.3s ease;
}
h1:hover {
    color: #c678dd;
}
h2 {
    color: #e06c74;
    margin-bottom: 20px;
    font-size: 2.5em;
    transition: color 0.3s ease;
}
h2:hover {
    color: #d4be98;
}
table {
    width: 100%;
    border-collapse: collapse;
    margin-bottom: 40px;
    box-shadow: 0 2px 3px rgba(0,0,0,0.1);
    transition: box-shadow 0.3s ease;
}
table:hover {
    box-shadow: 0 4px 6px rgba(0,0,0,0.2);
}
th, td {
    padding: 15px;
    border-bottom: 1px solid #3e4452;
    vertical-align: top;
    transition: background-color 0.3s ease;
}
th {
    background-color: #444;
    color: #ffffff;
}
tr:nth-child(even){
    background-color: #21252b;
}
tr:hover{
    background-color: #323642;
    cursor: pointer;
}
pre {
    background-color: #2c3e50;
    color: #abb2bf;
    padding: 15px;
    border-radius: 8px;
    white-space: pre-wrap;
    word-wrap: break-word;
    transition: background-color 0.3s ease;
}
pre:hover {
    background-color: #3e4452;
}
.line {
    font-size: 14px;
    color: #81a1c1;
    transition: color 0.3s ease;
}
.line:hover {
    color: #d19a66;
}
.analysis {
    margin-top: 10px;
    border-left: 4px solid #e06c74;
    padding-left: 10px;
    transition: border-color 0.3s ease;
}
.analysis:hover {
    border-color: #d4be98;
}
.analysis-header {
    font-size: 18px;
    color: #e06c74;
    margin-bottom: 10px;
    transition: color 0.3s ease;
}
.analysis-header:hover {
    color: #c678dd;
}
.analysis-desc {
    font-size: 14px;
    color: #b69828;
    margin-bottom: 8px;
    transition: color 0.3s ease;
}
.analysis-desc:hover {
    color: #56b6c2;
}
.analysis-payload {
    font-size: 14px;
    color: #98c379;
    margin-bottom: 8px;
    transition: color 0.3s ease;
}
.analysis-payload:hover {
    color: #98c379;
}
.analysis-code {
    font-size: 14px;
    color: #d19a66;
    transition: color 0.3s ease;
}
.analysis-code:hover {
    color: #c678dd;
}
</style>
</head>
<body>
<h1>Vulnerability Report</h1>
{{range .}}
<h2>{{.FileName}}</h2>
<table>
<tr>
    <th>File Path</th>
    <th>Line</th>
    <th>Code</th>
    <th>Analysis</th>
</tr>
{{range .Vulnerabilities}}
<tr>
    <td>{{.FilePath}}</td>
    <td>{{.Line}}</td>
    <td><pre class="analysis-code">{{.Code}}</pre></td>
    <td>
        {{if .VulnType}}
        <div class="analysis-header">AI Analysis Result:</div>
        <div class="analysis-desc"><strong>Vulnerability Type:</strong> {{.VulnType}}</div>
        <div class="analysis-desc"><strong>Severity:</strong> {{.Severity}}</div>
        <div class="analysis-desc"><strong>Reasons:</strong> {{.Description}}</div>
        <div class="analysis-payload"><strong>Payload:</strong> <pre>{{.Payload | html}}</pre></div>
        {{else}}
        <div>大概率没有漏洞</div>
        {{end}}
    </td>
</tr>
{{end}}
</table>
{{end}}
</body>
</html>
`))

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}
