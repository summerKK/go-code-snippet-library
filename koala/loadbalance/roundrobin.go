package loadbalance

import (
	"context"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
)

type Roundrobin struct {
	name  string
	index int
}

func NewRoundrobin() *Roundrobin {
	return &Roundrobin{
		name:  "roundrobin",
		index: 0,
	}
}

func (r *Roundrobin) GetName() string {
	return r.name
}

func (r *Roundrobin) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	if len(nodes) == 0 {
		err = ErrNoNodes
		return
	}

	r.index = (r.index + 1) % len(nodes)
	node = nodes[r.index]
	return
}
