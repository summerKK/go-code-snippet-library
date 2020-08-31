package logic

import (
	"strings"

	"github.com/summerKK/go-code-snippet-library/chatroom/global"
)

func FilterSensitiveWord(word string) string {
	for _, sensitiveWord := range global.SensitiveWords {
		word = strings.ReplaceAll(word, sensitiveWord, "**")
	}

	return word
}
