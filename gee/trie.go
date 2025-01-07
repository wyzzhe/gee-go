package gee

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

func (n *node) matchChild(part string) *node {
	// 遍历切片 index[0,1...], value
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 将路由模式（pattern）插入到前缀树中
// pattern /p/:lang/doc
// parts ["p", ":lang", "doc"]
func (n *node) insert(pattern string, parts []string, height int) {
	// 递归出口条件：处理完所有分段就返回
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 获取当前part
	part := parts[height]
	// 根据part查找单个子节点
	child := n.matchChild(part)
	// 如果没有找到匹配的子节点，则创建一个新的子节点
	if child == nil {
		// 创建子节点
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		// 添加子节点至当前节点 n 的 children 列表中
		n.children = append(n.children, child)
	}
	// 递归插入节点
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// 递归终止条件：
	// 1. 如果当前处理的索引 height 等于 parts 的长度，说明已经处理完所有分段
	// 2. 如果当前节点的 part 以 * 开头，说明这是一个通配符节点，匹配任意内容。
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	// 查找子节点列表
	children := n.matchChildren(part)

	// 递归查找
	for _, child := range children {
		// 遍历所有匹配的子节点，递归调用 search 方法，处理下一个分段（height+1）
		result := child.search(parts, height+1)
		// 找到匹配的节点后返回该节点
		if result != nil {
			return result
		}
	}

	// 未查找到
	return nil
}
