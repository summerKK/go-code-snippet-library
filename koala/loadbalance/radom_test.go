package loadbalance

import (
	"context"
	"fmt"
	"github.com/summerKK/go-code-snippet-library/koala/registry"
	"testing"
)

func TestRandom_Select(t *testing.T) {
	random := &Random{}
	var nodes []*registry.Node
	nodeWeight := [3]int{50, 100, 150}
	for i := 0; i < 4; i++ {
		n := &registry.Node{
			Id:     i,
			Ip:     fmt.Sprintf("127.0.0.%d", i),
			Port:   3306,
			Weigth: nodeWeight[i%len(nodeWeight)],
		}
		nodes = append(nodes, n)
		fmt.Printf("node:%s,weight:%d\n", n.Ip, n.Weigth)
	}

	fmt.Println()

	nodeCount := make(map[string]int)
	for i := 0; i < 100; i++ {
		node, err := random.Select(context.TODO(), nodes)
		if err != nil {
			fmt.Println(err)
		}
		nodeCount[node.Ip]++
	}

	for k, count := range nodeCount {
		fmt.Printf("node:%v,count:%d\n", k, count)
	}
}
