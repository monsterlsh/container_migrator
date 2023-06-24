/*
 */
package dirty

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/JBinin/container-migrator/util"
)

func DifferChkpdir() error {
	chkpList, err := util.GetCheckpointDir(checkpointPath)
	if err != nil {
		log.Panic("get chkpList failed ")
	}
	//oldDir, _ = os.Getwd()
	os.Chdir(csvPath)
	oldDir, _ := os.Getwd()
	log.Println("in DifferChkpdir() ==> ", oldDir)
	//n := len(chkpList) - 1
	var iter string
	iterNext := filepath.Join(checkpointPath, chkpList[0])
	var iterPagemap string
	var iterNextPagemap string
	var pagemap_iter Pagemap
	var pagemap_iterNext Pagemap

	var dirtyPagesEntry []CsvEntry
	var tmp []Entry
	//bitmap, err := GetAllPages()
	for k, chkp := range chkpList {
		if k == 0 {
			continue
		}
		iter = iterNext
		iterNext = filepath.Join(checkpointPath, chkp)
		set_iter, err := util.GetImageFile(iter)
		if err != nil {
			return err
		}
		list_iter := set_iter.List()
		sort.Strings(list_iter)
		set_iterNext, err := util.GetImageFile(iterNext)
		if err != nil {
			return err
		}
		list_iterNext := set_iterNext.List()
		sort.Strings(list_iterNext)
		log.Printf("diff iter=%s and iterNext=%s", iter, iterNext)
		for i, j := 0, 0; i < len(list_iter); i, j = i+1, j+1 {
			if list_iter[i] != list_iterNext[i] {
				continue
			}
			iterPagemap = filepath.Join(iter, list_iter[i])
			iterNextPagemap = filepath.Join(iterNext, list_iter[i])
			log.Printf("\t\\__ diff %s ", list_iter[i])
			pagemap_iter, err = readPagemapJSON(iterPagemap)
			if err != nil {
				log.Fatalln(err)
			}
			pagemap_iterNext, err = readPagemapJSON(iterNextPagemap)
			if err != nil {
				log.Fatalln(err)
			}
			tmp, err = ComparePagemaps(pagemap_iter, pagemap_iterNext)
			if err != nil {
				log.Fatalln(err)
			}
			csvtmp := CsvEntry{iteration: i, Entries: tmp}
			dirtyPagesEntry = append(dirtyPagesEntry, csvtmp)
		}

	}
	// dirtyPages := JsonEntry{CsvEntry: dirtyPagesEntry}
	// jsonData, err := json.Marshal(dirtyPages)
	// if err != nil {
	// 	log.Fatalf("Error marshalling JSON: %v", err)
	// }
	// fmt.Println(string(jsonData))
	ConvertJSONToCSV(dirtyPagesEntry, dirtypagesPath)
	log.Println("in DifferChkpdir() ==> done")

	return nil
}
