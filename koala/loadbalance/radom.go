package loadbalance

import (
	"context"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"math/rand"
)

type Random struct {
}

func (r *Random) Name() string {
	return "radom"
}

func (r *Random) Select(ctx context.Context, nodes []*registry.Node) (selectedNode *registry.Node, err error) {
	if len(nodes) == 0 {
		err = ErrNoNodes
		return
	}

	countWeight := 0
	for _, node := range nodes {
		// 设置默认权重
		if node.Weigth == 0 {
			node.Weigth = 1
		}
		countWeight += node.Weigth
	}

	randWeight := rand.Intn(countWeight)
	currentNodeIndex := -1
	for i, node := range nodes {
		randWeight -= node.Weigth
		if randWeight < 0 {
			currentNodeIndex = i
			break
		}
	}
	if currentNodeIndex == -1 {
		err = ErrNoNodes
	}

	selectedNode = nodes[currentNodeIndex]
	return
}
