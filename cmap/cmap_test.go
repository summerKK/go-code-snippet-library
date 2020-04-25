package cmap

import (
	"fmt"
	"log"
	"testing"
)

var cmp IConcurrentMap

func init() {
	var err error
	cmp, err = NewConcurrentMap(100, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func TestConcurrentMap_Put(t *testing.T) {
	_, err := cmp.Put("summer", "陈思贝")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	v := cmp.Get("summer")
	fmt.Printf("%s", v)

}
