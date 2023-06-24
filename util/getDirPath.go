package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func fillBool(v []bool, value bool) {
	for i := 0; i < len(v); i++ {
		v[i] = value
	}
	//return v
}
func swap(s []string, i int, num int) {
	ibyte := []byte(s[i])
	numbyte := []byte(s[num])
	ibyte, numbyte = numbyte, ibyte
	s[i] = string(ibyte)
	s[num] = string(numbyte)
}
func printArray(subPaths []string) {
	for i, e := range subPaths {
		if i != len(subPaths)-1 {
			fmt.Printf("%s,", e)
		} else {
			fmt.Printf("%s\n", e)
		}
	}
}
func sortSubPath(s []string) []string {
	n := 0
	len_s := len(s)
	var e string
	var err error
	var tmp int
	v := make([]bool, len_s)
	fillBool(v, false)
	i := 0
	for i < len_s {
		if v[i] {
			i++
			continue
		}
		e = s[i]
		//log.Printf("now at i=%d ,e =%s", i, e)
		num := -1
		n = len(e)
		for idx := len(e) - 1; idx >= 0; idx-- {
			tmp, err = strconv.Atoi(e[idx:n])
			if err != nil {
				break
			}
			num = tmp
		}
		//log.Printf("now at i=%d ,e =%s , num=%d,snum=%s", i, e, num, s[num])
		if num == -1 {
			swap(s, i, len_s-1)

			v[len_s-1] = true
			len_s--
			continue
		}
		if i == num {
			v[i] = true
			i++
			continue
		}
		swap(s, i, num)
		v[num] = true
	}
	return s[:len_s]
}
func GetCheckpointDir(root string) ([]string, error) {
	subPaths := make([]string, 0)
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
			//log.Printf("name = %s\n", info.Name())
			subPaths = append(subPaths, info.Name())
		}
		return nil
	})
	if err != nil {
		return subPaths, err
	}
	//printArray(subPaths)
	subPaths = sortSubPath(subPaths)
	//printArray(subPaths)
	return subPaths, nil
}
