package basetest

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	sl := []string{"mumbai", "london", "tokyo", "seattle"}
	sort.Strings(sl)
	fmt.Println(sl)
}
