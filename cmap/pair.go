package cmap

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"unsafe"
)

type IPair interface {
	IlinkedPair
	// 返回键的值
	Key() string
	// 返回键的散列值
	Hash() uint64
	// 返回元素的值
	Element() interface{}
	// 设置元素的值
	SetElement(element interface{}) error
	// 返回一个副本
	Copy() IPair
	String() string
}

type IlinkedPair interface {
	// 若返回nil,代表已在单链表的末尾
	Next() IPair
	SetNext(nextPair IPair) error
}

type pair struct {
	key     string
	hash    uint64
	element unsafe.Pointer
	next    unsafe.Pointer
}

func newPair(key string, element interface{}) (*pair, error) {
	if element == nil {
		return nil, newIllegalParameterError("element is nill")
	}
	newPair := &pair{
		key:  key,
		hash: hash(key),
	}
	// element
	newPair.element = unsafe.Pointer(&element)

	return newPair, nil
}

func (p *pair) Next() IPair {
	pointer := atomic.LoadPointer(&p.next)
	if pointer == nil {
		return nil
	}

	return (*pair)(pointer)
}

func (p *pair) SetNext(nextPair IPair) error {
	if nextPair == nil {
		return newIllegalParameterError("nextPair is nil")
	}
	pp, ok := nextPair.(*pair)
	if !ok {
		return newIllegalPairTypeError(nextPair)
	}
	atomic.StorePointer(&p.next, unsafe.Pointer(pp))

	return nil
}

func (p *pair) Key() string {
	return p.key
}

func (p *pair) Hash() uint64 {
	return p.hash
}

func (p *pair) Element() interface{} {
	// 原子操作
	pointer := atomic.LoadPointer(&p.element)
	if pointer == nil {
		return nil
	}
	return *(*interface{})(pointer)
}

func (p *pair) SetElement(element interface{}) error {
	if element == nil {
		return newIllegalParameterError("element is nil")
	}
	atomic.StorePointer(&p.element, unsafe.Pointer(&element))

	return nil
}

func (p *pair) Copy() IPair {
	pCopy, _ := newPair(p.Key(), p.Element())
	return pCopy
}

func (p *pair) String() string {
	return p.genString(false)
}

// genString 用于生成并返回当前键-元素对的字符串形式。
func (p *pair) genString(nextDetail bool) string {
	var buf bytes.Buffer
	buf.WriteString("pair{key:")
	buf.WriteString(p.Key())
	buf.WriteString(", hash:")
	buf.WriteString(fmt.Sprintf("%s", p.Hash()))
	buf.WriteString(", element:")
	buf.WriteString(fmt.Sprintf("%+v", p.Element()))
	if nextDetail {
		buf.WriteString(", next:")
		if next := p.Next(); next != nil {
			if npp, ok := next.(*pair); ok {
				buf.WriteString(npp.genString(nextDetail))
			} else {
				buf.WriteString("<ignore>")
			}
		}
	} else {
		buf.WriteString(", nextKey:")
		if next := p.Next(); next != nil {
			buf.WriteString(next.Key())
		}
	}
	buf.WriteString("}")
	return buf.String()
}
