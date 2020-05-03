package buffer

import (
	"errors"
)

var (
	// 缓冲池关闭
	ErrClosedBufPool = errors.New("closed buffer pool")
	// 缓冲器关闭
	ErrClosedBuf = errors.New("closed buffer")
)
