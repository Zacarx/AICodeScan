package Utils

import (
	"AICodeScan/CommonVul/Rce"
	"AICodeScan/CommonVul/Upload"
	"AICodeScan/PHP-Code/FileRead"
	"AICodeScan/PHP-Code/Include"
	"AICodeScan/PHP-Code/PHPSql"
	"AICodeScan/PHP-Code/SSRF"
	"AICodeScan/PHP-Code/Unserialize"
	"os"
	"path/filepath"
	"strings"
)

func PHP_Codeing() {
	//StartTime = time.Now()

	// 所有要执行的扫描函数
	scanFuncs := []func(string){
		Upload.PHPUpload_check,
		Rce.PHPRce,
		PHPSql.Sqlcheck,
		FileRead.Read,
		Unserialize.Unserialize,
		SSRF.PHP_SSRF,
		Include.Include,
	}

	//var wg sync.WaitGroup
	//wg.Add(len(scanFuncs)) // 根据方法数量动态调整 goroutine 数量
	//progressBar = pb.New(len(scanFuncs)).SetRefreshRate(time.Millisecond * 100).Start()
	// 启动 goroutine 来执行扫描任务
	for _, scanFunc := range scanFuncs {
		scanDirectory(scanFunc, *Dir)
	}

	//wg.Wait()
	//progressBar.Finish()
	// 清理空文件
	root := "./" // 设置要检查的目录
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			if info.Size() == 0 {
				os.Remove(path)
			}
		}
		return nil
	})
}
