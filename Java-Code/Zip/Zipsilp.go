package Zip

import (
	"AICodeScan/FindFile"
	"fmt"
)

func Zipsilp(dir string) {
	fmt.Println("Zipsilp分析开始")
	FindFile.FindFileByJava(dir, "zip.txt", []string{"zipEntry.getName(", "ZipUtil.unpack(", "ZipUtil.unzip(", "entry.getName()", "AntZipUtils.unzip(", "zip.getEntries()"})
	fmt.Println("Zipsilp分析完成")
}
