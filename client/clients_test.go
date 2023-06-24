package client

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func GetWorkingDirPath() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	s2 := filepath.Join(dir, "../", "workloads/redis")

	return s2
}
func TestGetWorkingDirPath(t *testing.T) {
	othersPath := GetWorkingDirPath()

	log.Println(othersPath)
}
func TestGetAllFile(t *testing.T) {
	var dir = GetWorkingDirPath()
	log.Printf("list all of %s\n", dir)

	s, err := GetAllFile(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(s)
}
func TestConnection(t *testing.T) {
	var containerID = "myredis"
	var destIP = "127.0.0.1"

	var othersPath = GetWorkingDirPath()
	basePath := path.Join(client_path, containerID)

	//buf := [512]byte{}
	Connection(containerID, destIP, othersPath, basePath)
}
func TestPreDump(t *testing.T) {
	var containerID = "myredis"
	oldDir, _ := os.Getwd()
	// if err := os.Chdir("/opt/test_go_migrator/workloads/redis"); err != nil {
	// 	log.Println("Failed to change the work directory")
	// }
	// path, _ := os.Getwd()
	// log.Println("the workspace is", path)
	// if _, err := exec.Command("runc", []string{"run", containerID, ">", "/dev/null", "&"}...).Output(); err != nil {
	// 	log.Println(err.Error())

	// }
	// if err := os.Chdir("/opt/migrator/client"); err != nil {
	// 	log.Println("Failed to change the work directory")
	// }
	path, _ := os.Getwd()
	log.Println("the workspace is", path)
	defer os.Chdir(oldDir)
	// cmd := exec.Command("sleep", "10")
	// err := cmd.Run() //执行到此处时会阻塞等待10秒
	// if err != nil {
	// 	log.Fatal(err)

	// }

	log.Println("--------------predump--------------")
	if _, err := preDump(containerID, 0); err != nil {
		log.Println("iteration pre dump failed ")

	} else {
		log.Println("OK!!")
	}

	// if _, err := exec.Command("runc", []string{"run", "kill", containerID}...).Output(); err != nil {
	// 	log.Println(err.Error())
	// }

}
func TestIterator(t *testing.T) {
	var containerID = "sys"
	var destIP = "127.0.0.1"

	destPath := path.Join(server_path, containerID)
	basePath := path.Join(client_path, containerID)
	if err := os.Chdir(basePath); err != nil {
		log.Println("Failed to change the work directory")
	} else {
		pwd, _ := os.Getwd()
		log.Println("current workspace:", pwd)
	}

	if _, err := iterator(containerID, basePath, destIP, destPath); err != nil {
		log.Println("Iterator transfer failed")

	}
}
func TestRunc(t *testing.T) {
	RunMyredis()
}
func TestRestore(t *testing.T) {
	containerID := "myredis"
	imagePath := path.Join(server_path, "checkpoint")
	args := []string{"restore", "-j", "-d", "--auto-dedup", "--tcp-established", "--image-path", imagePath, containerID}

	oldDir, _ := os.Getwd()
	os.Chdir(server_path)
	defer os.Chdir(oldDir)

	start := time.Now()
	if err := exec.Command("runc", args...).Run(); err != nil {
		log.Println("Failed to restore the contaier")
		log.Println(err.Error())
		return
	}
	elapsed := time.Since(start)
	log.Println("Restore time is ", elapsed)
}
func TestPreCopy(t *testing.T) {

	//RunMyredis()
	log.SetFlags(log.Ldate | log.Lshortfile)
	var containerID = "kafka"
	var destIP = "127.0.0.1"
	// get the dir of rootfs in client
	dir, _ := os.Getwd()
	othersPath := filepath.Join(dir, "../", "workloads", containerID)

	if err := PreCopy(containerID, destIP, othersPath); err != nil {
		log.Println("precopy failed")

	}
}
func TestGetDirty(t *testing.T) {
	log.SetFlags(log.Ldate | log.Lshortfile)
	var containerID = "mysql"
	basePath := path.Join(client_path, containerID)
	log.Println(basePath)
	// remove all the files in the basepath

	err := filepath.Walk(basePath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("遍历文件路径出错：%v\n", err)
			return nil
		}
		log.Println("delete the file:", filePath)
		if filePath != basePath {
			err = os.RemoveAll(filePath)
			if err != nil {
				log.Printf("delete wrong：%v\n", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Println("删除子文件失败:", err)
		return
	}

	log.Println("子文件删除成功")

	//get the dir of rootfs in client
	//time.Sleep(time.Second * 5)
	var index int
	for i := 0; i < 500; i += 1 {
		index = i
		if preTime, err := preDump(containerID, i); err != nil {
			log.Println("The ", index, "iteration pre dump failed ")
			break
		} else {
			preDumpPath := path.Join(basePath, "checkpoint"+strconv.Itoa(index))
			log.Println("\t\\__iteration()-> preDumpPath is", preDumpPath)
			log.Printf("\t\\__iterator() ==> in preDump() consuming %f (s)", preTime)
		}
		time.Sleep(time.Second)
	}
}
