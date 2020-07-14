package main

import (
	"log"

	"github.com/summerKK/go-code-snippet-library/word-conversion/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute error:%v", err)
	}
}
