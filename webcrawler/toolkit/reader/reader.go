package reader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

//  多重读取器
type MultipleReader struct {
	data []byte
}

func NewReader(reader io.Reader) (*MultipleReader, error) {
	var data []byte
	if reader != nil {
		readAll, err := ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("multiple reader:could't create new one:%s", err)
		}
		data = readAll
	} else {
		data = []byte{}
	}

	return &MultipleReader{
		data: data,
	}, nil
}

func (r *MultipleReader) Reader() io.ReadCloser {
	// ioutil.NopCloser 实现一个不需要关闭的ReadCloser.
	return ioutil.NopCloser(bytes.NewReader(r.data))
}
