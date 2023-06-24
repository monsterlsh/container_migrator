package client

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// T the max expected time of downtime(s)
var T float64 = 1

const (
	client_path string = "/opt/container-migrator/client_repo"
	server_path string = "/opt/container-migrator/server_repo/targetv0.1"
)

type transferInfo struct {
	index        int
	data         float64
	preTime      float64
	transferTime float64
}

var Info []transferInfo

func PrintInfo() {
	log.Println("---------------------PrintInfo--------------------------------------")
	log.Println("index\t", "data-size(KB)\t\t", "pre-time(s)\t", "transfer-time(s)\t")
	for _, f := range Info {
		log.Println(f.index, "\t\t\t", f.data, "\t\t\t", f.preTime, "\t\t", f.transferTime)
	}
	log.Println("--------------------------------------------------------------------")
}

func preDump(containerID string, index int) (preTime float64, err error) {
	checkpointDestPath := path.Join(client_path, containerID, "checkpoint"+strconv.Itoa(index))
	log.Println("in preDump() chekpoint Path :", checkpointDestPath)

	start := time.Now()
	args := []string{
		"checkpoint",
		"--pre-dump",
		//"--auto-dedup",
		"--tcp-established",
		"--image-path",
		checkpointDestPath,
	}
	if index != 0 {
		//parentPath := path.Join(client_path, containerID, "checkpoint"+strconv.Itoa(index-1))
		parentPath := path.Join("..", "checkpoint"+strconv.Itoa(index-1))
		log.Println("predump()-> parentPath is ", parentPath)
		args = append(args, "--parent-path", "../checkpoint"+strconv.Itoa(index-1))
		//args = append(args, "--parent-path", parentPath)
	}
	args = append(args, containerID)
	if output, err := exec.Command("runc", args...).Output(); err != nil {
		log.Println(output)
		log.Println(err.Error())

		return 0, err
	}
	elapsed := time.Since(start)
	log.Println("The pre-dump index is ", index, " . The pre-dump time is ", elapsed)
	return elapsed.Seconds(), nil
}

func dump(containerID string, index int) (dumpTime float64, err error) {

	start := time.Now()
	args := []string{
		"checkpoint",
		"--auto-dedup",
		"--tcp-established",
		"--image-path",
		"checkpoint",
		"--parent-path",
		"../checkpoint" + strconv.Itoa(index),
		containerID,
	}
	if _, err := exec.Command("runc", args...).Output(); err != nil {
		log.Println(err.Error())
		return 0, err
	}
	elapsed := time.Since(start)
	return elapsed.Seconds(), nil
}

func transfer(sourcePath string, destIP string, destPath string, otherOpts []string) (transferTime float64, size int, err error) {
	if output, err := exec.Command("du", "-s", sourcePath).Output(); err != nil {
		log.Printf("du -s %s failed ", sourcePath)
		return 0, 0, err
	} else {
		size, _ = strconv.Atoi(strings.Split(string(output), "\t")[0])
		log.Println("Transfer size: ", size, " KB")
	}
	dest := destIP + ":" + destPath
	rsyncOpts := []string{"-aqz", "--bwlimit=125000", sourcePath, dest}
	//if otherOpts != nil {
	//	//rsyncOpts = append(otherOpts, rsyncOpts...)
	//}
	start := time.Now()
	if _, err := exec.Command("rsync", rsyncOpts...).Output(); err != nil {
		log.Println(err.Error())
		return 0, size, err
	}
	elapsed := time.Since(start)
	return elapsed.Seconds(), size, nil
}

/*
*
*
*
 */
func iterator(containerID string, basePath string, destIP string, destPath string) (int, error) {
	var index int
	D := 1e5
	N := 1.25e5
	S := T * (D * N / (2*N + D)) * 1024 / 10000
	log.Println("-----------------------------------")
	log.Println("Disk IO : ", D, " KB/s")
	log.Println("Net speed: ", N, " KB/s")
	log.Println("Expect memory size: ", S, "KB")
	log.Println("-----------------------------------")
	//var preDumpPath string
	for i := 0; i < 100; i += 1 {
		index = i
		if preTime, err := preDump(containerID, i); err != nil {
			log.Println("The ", index, "iteration pre dump failed ")
			return index, err
		} else {
			preDumpPath := path.Join(basePath, "checkpoint"+strconv.Itoa(index))
			//otherOpts := []string{"--remove-source-files"}
			log.Println("\t\\__iteration()-> preDumpPath is", preDumpPath, " destPath is ", destPath)
			log.Printf("\t\\__iterator() ==> in preDump() consuming %f (s)", preTime)
			// if transferTime, size, err := transfer(preDumpPath, destIP, destPath, otherOpts); err != nil {
			// 	log.Println("The ", index, "iteration transfer pre data failed")
			// 	return index, err
			// } else {
			// 	Info = append(Info, transferInfo{
			// 		index:        index,
			// 		data:         float64(size),
			// 		preTime:      preTime,
			// 		transferTime: transferTime,
			// 	})
			// 	if float64(size) < S {
			// 		log.Println("shut down the limit size", S)
			// 		break
			// 	}
			// }
		}
	}
	return index, nil
}
func GetAllFile(pathname string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		log.Printf("\tnow is show %s", fi.Name())
		if fi.IsDir() {
			fullDir := filepath.Join(pathname, fi.Name())
			s = append(s, fullDir)
		} else {
			fullName := filepath.Join(pathname, fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}

func syncReadOnly(destPath string, destIP string, othersPath string) error {
	log.Printf("transfer %s to %s:%s\n", othersPath, destIP, destPath)
	s, err := GetAllFile(othersPath)
	if err != nil {
		log.Printf("ERROR in get subpath of %s", othersPath)
	}
	for _, subfile := range s {

		if subfile == "data" {
			continue
		}
		if transferTime, size, err := transfer(subfile, destIP, destPath, nil); err != nil {
			log.Printf("Failed to sync the %s", subfile)
			return err
		} else {
			log.Printf("-----------------%s------------------\n", subfile)
			log.Println("data-size(KB) : ", size, "\t", "transfer time(s): ", transferTime)
			log.Println("----------------------------------------------")
		}
		// if transferTime, size, err := transfer(path.Join(othersPath, "rootfs"), destIP, destPath, nil); err != nil {
		// 	log.Println("Failed to sync the rootfs")
		// 	return err
		// } else {
		// 	log.Println("--------------------rootfs--------------------")
		// 	log.Println("data-size(KB) : ", size, "\t", "transfer time(s): ", transferTime)
		// 	log.Println("----------------------------------------------")
		// }
	}
	return nil
}

func syncVolume(destPath string, destIP string, othersPath string) error {
	if transferTime, size, err := transfer(path.Join(othersPath, "data"), destIP, destPath, nil); err != nil {
		log.Println("Failed to sync the volume")
		return err
	} else {
		log.Println("----------------volume----------------------")
		log.Println("data-size(KB) : ", size, "\t", "transfer time(s): ", transferTime)
		log.Println("--------------------------------------------")
	}
	return nil
}
func GetDestPath(conn net.Conn, buf [512]byte) (string, error) {
	var destPath string
	if n, err := conn.Read(buf[:]); err != nil {

		return destPath, err
	} else {
		destPath = string(buf[:n])
	}

	return destPath, nil
}

func Connection(containerID string, destIP string, othersPath string, basePath string) (net.Conn, error) {
	var conn net.Conn
	var conErr error
	conn, conErr = net.Dial("tcp", destIP+":8001")
	//defer conn.Close()
	if conErr != nil {
		log.Println("Tcp connect failed")
		return nil, conErr
	}

	if err := os.RemoveAll(basePath); err != nil {
		log.Println("Remove ", basePath, " failed")
		return nil, err
	}
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		log.Println("Mkdir ", basePath, " failed")
		return nil, err
	}

	if _, err := conn.Write([]byte(containerID)); err != nil {
		log.Println("Send container id or get DestPath failed")
		return nil, err
	}

	return conn, nil
}
func PreCopy(containerID string, destIP string, othersPath string) error {
	defer PrintInfo()
	//currentTime := time.Now()
	buf := [512]byte{}
	basePath := path.Join(client_path, containerID)
	// if f, err := util.PathExists(basePath); !f {
	// 	log.Printf("\\__in client.go PreCopy() ==> Error opening %s file: %v\n", basePath, err)

	// } else {
	// 	basePath = basePath + currentTime.Format("2006-01-01 15:04:05")
	// 	for err := os.Mkdir(basePath, 7777); err != nil; {
	// 		log.Printf("fail mkdir %s", basePath)
	// 		os.RemoveAll(basePath)
	// 		//return err
	// 	}
	// }
	conn, err := Connection(containerID, destIP, othersPath, basePath)
	//connect to server
	if err != nil {
		return err
	}

	//get time and path
	totalStart := time.Now()
	oldDir, _ := os.Getwd()

	//get destPath
	destPath, err := GetDestPath(conn, buf)
	if err != nil {
		log.Println("Get DestPath failed")
	}

	// sync container rootfs from client to server
	if err := syncReadOnly(destPath, destIP, othersPath); err != nil {
		log.Println("Sync readonly dir failed")
		return err
	}
	// change workspace to path 'workloads/redis'
	if err := os.Chdir(basePath); err != nil {
		log.Println("Failed to change the work directory")
		return err
	}
	defer os.Chdir(oldDir)
	//checkpointDestPath := path.Join(server_path, containerID)
	if index, err := iterator(containerID, basePath, destIP, destPath); err != nil {
		log.Println("Iterator transfer failed")
		return err
	} else {
		start := time.Now()
		if dumpTime, err := dump(containerID, index); err != nil {
			log.Println("Dump data failed")
			return err
		} else {
			dumpPath := path.Join(basePath, "checkpoint")
			otherOpts := []string{"--remove-source-files"}
			if transferTime, size, err := transfer(dumpPath, destIP, destPath, otherOpts); err != nil {
				log.Println("Transfer dump data failed")
				return err
			} else {
				if err := syncVolume(destPath, destIP, othersPath); err != nil {
					log.Println("Failed to sync the volume")
				}
				log.Println("---------------------dump------------------------")
				log.Println("dumpTime(s)\t", "data-size(KB)\t", "transfer time(s)")
				log.Println(dumpTime, "\t", size, "\t", transferTime, "\t")
				log.Println("-------------------------------------------------")
			}
		}
		if _, err := conn.Write([]byte("restore:" + containerID)); err != nil {
			log.Println("Send restore cmd failed")
			return err
		}
		if n, err := conn.Read(buf[:]); err != nil {
			log.Println("Waiting for restore container in another machine")
			return err
		} else {
			if string(buf[:n]) == "started" {
				elapsed := time.Since(start)
				log.Println("The downtime is ", elapsed)
			} else {
				log.Println("Restore error in remote machine")
				return errors.New("Restore failed ")
			}
		}
	}
	totalElapsed := time.Since(totalStart)
	log.Println("The total migration time is ", totalElapsed)
	return nil
}

func PostCopy(containerID string, destination string) error {
	// todo
	return nil
}
