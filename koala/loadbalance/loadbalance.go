package loadbalance

import (
	"context"
	"errors"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
)

var (
	ErrNoNodes = errors.New("没有节点")
)

type ILoadbalance interface {
	GetName() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}
