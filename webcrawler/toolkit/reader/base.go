package reader

import "io"

type IReader interface {
	Reader() io.ReadCloser
}
