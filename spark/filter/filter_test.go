package filter

import "testing"

func init() {
	err := Init("../data/filter.dat")
	if err != nil {
		panic(err)
	}
}

func TestFilter(t *testing.T) {
	text := "hello,world偷情"
	rText, replace := Filter(text, "***")
	if !replace || rText != "hello,world***" {
		t.Errorf("filter failed")
		return
	}
}
