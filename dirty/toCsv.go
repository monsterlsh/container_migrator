package dirty

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

// bitmap transfer to csv
// addr:[t1,...tn]
func ConvertMapoCSV(source map[string][timelens]int, destination string) error {
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// header := []string{"addr", "t", "nr_pages", "flag"}
	// if err := writer.Write(header); err != nil {
	// 	return err
	// }
	keys := make([]string, 0, len(source))
	for k := range source {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var v [timelens]int

	for _, k := range keys {
		v = source[k]
		var csvRow []string
		csvRow = append(csvRow, k)
		for _, e := range v {
			csvRow = append(csvRow, strconv.Itoa(e))
		}

		if err := writer.Write(csvRow); err != nil {
			return err
		}
	}
	return nil
}
func ConvertJSONToCSV(source []CsvEntry, destination string) error {

	// 3. Create a new file to store CSV data
	outputFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 4. Write the header of the CSV file and the successive rows by iterating through the JSON struct array
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"iteration", "vaddr", "nr_pages", "flag"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range source {

		entries := r.Entries
		for _, enetry := range entries {
			var csvRow []string
			csvRow = append(csvRow, strconv.Itoa(r.iteration), strconv.FormatInt(enetry.Vaddr, 16), strconv.Itoa(enetry.NrPages), strconv.Itoa(enetry.Flags))
			if err := writer.Write(csvRow); err != nil {
				return err
			}
		}
		//csvRow = append(csvRow, r.Vegetable, r.Fruit, fmt.Sprint(r.Rank))

	}
	return nil
}
