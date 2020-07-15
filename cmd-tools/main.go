package main

import (
	"log"

	"github.com/summerKK/go-code-snippet-library/cmd-tools/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute error:%v", err)
	}
}
