package connect

import (
	"context"

	"github.com/go-kratos/kratos/v2/selector"
)

// SelectNodeFilter 根据节点数据自定义过滤器
func SelectNodeFilter(filter func(node selector.Node) bool) selector.NodeFilter {
	return func(ctx context.Context, nodes []selector.Node) []selector.Node {
		newNodes := make([]selector.Node, 0, len(nodes))
		for _, node := range nodes {
			if filter(node) {
				newNodes = append(newNodes, node)
			}
		}
		return newNodes
	}
}
