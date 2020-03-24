package util

import (
	"crypto/md5"
	"fmt"
)

func Slat(data []byte) (slatStr string) {
	sum := md5.Sum(data)
	return fmt.Sprintf("%x", sum)
}
