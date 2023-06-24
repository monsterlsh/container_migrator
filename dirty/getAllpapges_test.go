package dirty

import (
	"log"
	"os"
	"testing"
)

func TestGetAllPages(t *testing.T) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	bitmap, err := GetAllPages()
	if err != nil {
		log.Fatalln("in TestGetAllPages() ==> wrong get all pages")
	}
	//log.Println("bitmap:\n\t", bitmap)
	err = ConvertMapoCSV(bitmap, allpagesPath)
	if err != nil {
		log.Println("fiale save", err)
	} else {
		log.Println("in GetAllPages() ==> done")
	}
}

func TestGetNumOfSubFilePath(t *testing.T) {
	basePath := "/opt/container-migrator/client_repo/mysql"
	dirEntries, err := os.ReadDir(basePath)
	if err != nil {
		log.Println("无法读取文件夹:", err)
		return
	}
	subDirCount := 0
	for _, entry := range dirEntries {
		if entry.IsDir() {
			subDirCount++
		}
	}
	log.Println("子文件夹数量:", subDirCount)
}
