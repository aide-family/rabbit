package connect

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"
)

type NodeFilter func(node selector.Node) bool

func SelectNodeFilterOr(filters ...NodeFilter) selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if len(filters) == 0 {
			return nodes
		}
		newNodes := make([]selector.Node, 0, len(nodes))
		for _, node := range nodes {
			anyPass := false
			for _, filter := range filters {
				if anyPass = anyPass || filter(node); anyPass {
					break
				}
			}
			if anyPass {
				newNodes = append(newNodes, node)
			}
		}
		return newNodes
	}
}

func SelectNodeFilterAnd(filters ...NodeFilter) selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		if len(filters) == 0 {
			return nodes
		}
		newNodes := make([]selector.Node, 0, len(nodes))
		for _, node := range nodes {
			allPass := true
			for _, filter := range filters {
				if allPass = allPass && filter(node); !allPass {
					break
				}
			}
			if allPass {
				newNodes = append(newNodes, node)
			}
		}
		return newNodes
	}
}
