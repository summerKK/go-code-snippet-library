package trie

import (
	"bufio"
	"io"
	"log"
	"os"
	"testing"
)

func TestTrie_Check(t *testing.T) {
	checkStr := "hello,world发生关系hello,world屁股"
	trie := NewTrie()
	file, err := os.Open("./filter.dat")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Errorf("read file failed")
			return
		}
		err = trie.Add(str, nil)
		if err != nil {
			t.Errorf("Trie add failed")
			return
		}
	}
	check, hit := trie.Check(checkStr, "***")
	if !hit || check != "hello,world***hello,world***" {
		t.Errorf("Trie check failed")
		return
	}
	log.Printf("%v\n", check)
}
