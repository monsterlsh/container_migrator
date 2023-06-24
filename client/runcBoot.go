package client

import (
	"log"
	"os"
	"path"
	"regexp"
)

func RunMyredis() {
	paths, _ := os.Getwd()
	log.Println(paths)
	newpath := path.Join(paths, "..", "workloads/redis")
	log.Printf("new workspace %s\n", newpath)
	datapath := path.Join(newpath, "data")
	s, _ := GetAllFile(datapath)
	for _, subfile := range s {
		reg := regexp.MustCompile(`.*sh`)
		if len(reg.FindAllString(subfile, -1)) > 0 {
			log.Println("the subfile si :", subfile)

		} else {

			if err := os.Remove(subfile); err != nil {
				log.Panicf("delete %s fail", subfile)
			}
		}

	}
	// cmd := exec.Command("runc", "run", "myredis")
	// err := cmd.Run()
	// if err != nil {
	// 	log.Fatalf("cmd.Run() failed with %s\n", err)
	// }
}
