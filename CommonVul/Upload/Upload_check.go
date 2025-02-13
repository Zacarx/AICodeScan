package Upload

import (
	"AICodeScan/CommonVul/Rule"
	"AICodeScan/FindFile"
	"fmt"
)

func JavaUpload_check(dir string) {
	fmt.Println("上传分析开始")
	//FindFile.FindFileByJava(dir, "upload.txt", []string{"new File(", "MultipartFile", "upload", ".getOriginalFilename(", ".transferTo("})
	FindFile.FindFileByJava(dir, "upload.txt", Rule.JavaUploadRuleList)
	fmt.Println("上传分析完成")
}

func PHPUpload_check(dir string) {
	fmt.Println("上传分析开始")
	FindFile.FindFileByPHP(dir, "upload.txt", Rule.PHPUploadRuleList)
	fmt.Println("上传分析完成")
}
