package filter

import (
	"bufio"
	"github.com/summerKK/go-code-snippet-library/trie"
	"io"
	"os"
)

var trieFilter *trie.Trie

func Init(filename string) (err error) {
	trieFilter = trie.NewTrie()
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	var readString string
	for {
		readString, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil
			return
		}
		if err != nil {
			return
		}
		err = trieFilter.Add(readString, nil)
		if err != nil {
			return
		}
	}
}
