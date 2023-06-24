package dirty

import (
	"fmt"
	"os"
	"path"
	"testing"
)

func TestComparePagemaps(t *testing.T) {
	checkpoint0 := path.Join(checkpointPath, "checkpoint0", "pagemap-1.json")
	checkpoint1 := path.Join(checkpointPath, "checkpoint1", "pagemap-1.json")
	pagemap1, err := readPagemapJSON(checkpoint0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading pagemap-1.json from checkpoint0: %v\n", err)
		os.Exit(1)
	}

	pagemap2, err := readPagemapJSON(checkpoint1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading pagemap-1.json from checkpoint1: %v\n", err)
		os.Exit(1)
	}

	ComparePagemaps(pagemap1, pagemap2)
}
