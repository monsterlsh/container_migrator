package dirty

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/JBinin/container-migrator/util"
)

//const pageSize = 4096

type Pagemap struct {
	Magic   string  `json:"magic"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	PagesID int   `json:"pages_id,omitempty"`
	Vaddr   int64 `json:"vaddr,omitempty"`
	NrPages int   `json:"nr_pages,omitempty"`
	Flags   int   `json:"flags,omitempty"`
}

func decodePagemaps(filename string, jsonfile string) error {
	//file := filepath.Base(filename) + ".json"
	cmd := exec.Command("crit", "decode", "-i", filename, "-o", jsonfile)
	if out, err := cmd.Output(); err != nil { //获取输出对象，可以从该对象中读取输出结果
		log.Println("\t\t\t\\__in decodePagemaps ==> wrong in crit")
		return err

	} else {
		log.Println("\t\t\t\\__in decodePagemaps crit out ==> ", out)
	}

	return nil
}

func readPagemapJSON(filename string) (Pagemap, error) {
	var pagemapEntries Pagemap
	filenameall := path.Base(filename)
	filesuffix := path.Ext(filename)
	fileprefix := filenameall[0 : len(filenameall)-len(filesuffix)]
	jsonname := filepath.Join(path.Dir(filename), fileprefix+".json")
	log.Println("\t\t\\__in readPagemapJSON ==> filename =", filename)
	if f, err := util.PathExists(jsonname); !f {
		log.Printf("\t\t\\__in readPagemapJSON ==> Error opening JSON file: %v\n", err)
		err = decodePagemaps(filename, jsonname)
		if err != nil {
			log.Panicln("\t\t\\__in readPagemapJSON() ==> decode ", filename, "failed")
			return pagemapEntries, err
		}
	}
	jsonFile, err := ioutil.ReadFile(jsonname)
	if err != nil {
		log.Panicln("\t\t\\__in readPagemapJSON() ==> ioutil.ReadFile ", filename, "failed")
		return pagemapEntries, err

	}
	//log.Println(string(jsonFile))

	err = json.Unmarshal(jsonFile, &pagemapEntries)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return pagemapEntries, nil
}

func ComparePagemaps(pagemap1, pagemap2 Pagemap) ([]Entry, error) {

	dirtyPages := make([]Entry, 0)
	entry1 := pagemap1.Entries
	entry2 := pagemap2.Entries
	json1Map := make(map[int64]Entry)
	for _, entry := range entry1 {
		for i := 0; i < entry.NrPages; i++ {
			json1Map[entry.Vaddr+int64(i*pageSize)] = entry
		}

	}

	var tmp_addr int64
	var index int
	sums := 0
	for _, entry := range entry2 {
		for i := 0; i < entry.NrPages; i++ {
			sums++
			tmp_addr = entry.Vaddr + int64(i*pageSize)
			entry1, exists := json1Map[tmp_addr]
			if !exists {
				//log.Printf("new pages %x", tmp_addr)
			} else if exists && entry1.Flags != entry.Flags {
				//log.Printf("dirty pages %x", tmp_addr)
				dirtyPages = append(dirtyPages, entry)
				index++
				break
			}
		}
	}

	log.Println("\t\t\t\\__in ComparePagemaps ==> Dirty pages: ", index, "all pages:", sums)
	// for _, entry := range dirtyPages {
	// 	fmt.Printf("\t\tVaddr: %x, NrPages: %d, Flags: %d\n", entry.Vaddr, entry.NrPages, entry.Flags)
	// }

	return dirtyPages, nil
}
