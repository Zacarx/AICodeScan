package Utils

import (
	"AICodeScan/CommonVul/Rce"
	"AICodeScan/CommonVul/Upload"
	"AICodeScan/Java-Code/AMF"
	"AICodeScan/Java-Code/Auth_Bypass"
	"AICodeScan/Java-Code/El"
	"AICodeScan/Java-Code/Fastjson"
	"AICodeScan/Java-Code/Frame_Analysis"
	"AICodeScan/Java-Code/JDBC"
	"AICodeScan/Java-Code/JNDI"
	"AICodeScan/Java-Code/JS"
	"AICodeScan/Java-Code/JarStatic"
	"AICodeScan/Java-Code/JavaSrciptShell"
	"AICodeScan/Java-Code/Log4j"
	"AICodeScan/Java-Code/ReadObject"
	"AICodeScan/Java-Code/Reflect"
	"AICodeScan/Java-Code/SSTI/FreeMarker"
	"AICodeScan/Java-Code/Sql"
	"AICodeScan/Java-Code/Zip"
	"os"
	"path/filepath"
	"strings"
)

func Java_Codeing() {
	//StartTime = time.Now()
	// 所有要执行的扫描函数
	scanFuncs := []func(string){
		Frame_Analysis.FrameAnalysiser,
		Auth_Bypass.Auth,
		Zip.Zipsilp,
		JNDI.Jndi,
		Sql.Sqlcheck,
		Rce.JavaRce,
		Upload.JavaUpload_check,
		ReadObject.Readobjectcheck,
		El.Elcheck,
		Fastjson.Parsecheck,
		Reflect.ReflectCheck,
		Log4j.Log4j,
		AMF.AmfCheck,
		FreeMarker.FreeSsti,
		JDBC.FindJDBC,
		JavaSrciptShell.FindJavaSrciptShell,
		JarStatic.Jarstaticer,
		JS.Eval,
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

	// 处理web.xml
	Frame_Analysis.WebXmlScan(*Dir, []string{"*.htm", "*.do", "*.action", "exclude"})

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
