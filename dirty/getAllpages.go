package dirty

import (
	"log"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/JBinin/container-migrator/util"
)

func GetAllPages() (map[string][timelens]int, error) {

	bitmap := make(map[string][timelens]int)
	chkpList, err := util.GetCheckpointDir(checkpointPath)
	if err != nil {
		log.Panic("get chkpList failed ")
	}

	var iter string
	var iterPagemap string
	var pagemap_iter Pagemap
	for k, chkp := range chkpList {

		iter = filepath.Join(checkpointPath, chkp)

		set_iter, err := util.GetImageFile(iter)
		if err != nil {
			return bitmap, err
		}
		list_iter := set_iter.List()
		sort.Strings(list_iter)

		for i := 0; i < len(list_iter); i++ {

			iterPagemap = filepath.Join(iter, list_iter[i])

			log.Printf("\t\\__ diff %s ", list_iter[i])

			pagemap_iter, err = readPagemapJSON(iterPagemap)
			if err != nil {
				log.Fatalln(err)
			}

			for _, entry := range pagemap_iter.Entries {
				tmp := entry.Vaddr
				var address string
				for j := 0; j < entry.NrPages; j++ {
					tmp += int64(j * pageSize)
					address = strconv.FormatInt(tmp, 16)
					if v, isContain := bitmap[address]; !isContain {
						newv := [timelens]int{}
						newv[k] = 1
						bitmap[address] = newv
					} else {
						v[k] = 1
						bitmap[address] = v
					}
				}
			}
		}
		if k == timelens-1 {
			break
		}
	}

	return bitmap, nil
}
