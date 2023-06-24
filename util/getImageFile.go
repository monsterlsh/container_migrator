package util

import (
	"os"
	"path/filepath"
)

func GetImageFile(root string) (*Set, error) {
	set := NewSet()
	_, file := filepath.Split(root)
	//log.Printf("filename = %s", file)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == file {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".img" {
			return nil
		}
		//log.Println(filepath.Ext(info.Name()))
		if info.Name()[0:7] != "pagemap" {
			//log.Println(info.Name(), info.Name()[0:4], info.Name()[4:7], info.Name()[0:4] == "page", info.Name()[4:7] == "map")
			return nil
		}

		//log.Printf("name = %s\n", info.Name())
		set.Add(info.Name())

		return nil
	})

	return set, err
}
