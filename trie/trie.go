package trie

import "strings"

type node struct {
	// utf8字符,兼容中文
	char rune
	// 父节点
	parent *node
	data   interface{}
	// 深度
	depth int
	// 叶子节点标识 true 代表叶子节点
	term bool
	// 子节点集合
	childs map[rune]*node
}

type trie struct {
	// 根节点
	root *node
	// 叶子节点总数
	size int
}

func newNode() *node {
	return &node{
		childs: make(map[rune]*node, 4),
	}
}

func NewTrie() *trie {
	return &trie{
		root: newNode(),
		size: 0,
	}
}

func (t *trie) Add(key string, data interface{}) (err error) {
	char := []rune(strings.TrimSpace(key))
	node := t.root
	for _, r := range char {
		// 查找当前node节点的childs是否存在r对应的节点
		child, ok := node.childs[r]
		if !ok {
			// 创建新节点
			child = newNode()
			child.depth += 1
			child.char = r
			child.parent = node
			// 把创建的新节点加入到当前node节点的childs
			node.childs[r] = child
		}
		// 当前node节点的childs存在r对应的节点,在 node.childs[r].childs节点继续查找
		node = child
	}

	// 叶子节点标识
	node.term = true
	// 叶子节点总数+1
	t.size += 1
	node.data = data
	return
}

// 查找节点
func (t *trie) findNode(text string) (result *node) {
	char := []rune(text)
	node := t.root
	for _, r := range char {
		child, ok := node.childs[r]
		if !ok {
			return
		}
		node = child
	}
	result = node
	return
}

// 给一个节点,然后手机该节点的所有子节点(包含自己)
func (t *trie) collectionNode(n *node) (result []*node) {
	if n == nil {
		return
	}

	// 当前节点就是一个叶子节点
	if n.term {
		result = append(result, n)
		return
	}

	queue := []*node{n}
	for i := 0; i < len(queue); i++ {
		// 遍历到叶子节点
		if queue[i].term {
			result = append(result, queue[i])
			// 这里不能return,因为节点可能有多个分支
			continue
		}
		// 把节点的子节点全部加入到遍历数组
		for _, child := range queue[i].childs {
			queue = append(queue, child)
		}
	}
	return
}

// 根据前缀查找子节点
func (t *trie) prefixSearch(key string) (result []*node) {
	node := t.findNode(key)
	if node == nil {
		return
	}
	result = t.collectionNode(node)
	return
}

func (t *trie) Check(text string, replace string) (rtext string, hit bool) {
	chars := []rune(strings.TrimSpace(text))
	if t.root == nil {
		return
	}
	var left []rune
	node := t.root
	start := 0
	for index, char := range chars {
		child, ok := node.childs[char]
		if !ok {
			left = append(left, chars[start:index+1]...)
			start = index + 1
			node = t.root
			continue
		}
		node = child
		if node.term {
			left = append(left, []rune(replace)...)
			node = t.root
			start = index + 1
			hit = true
		}
	}
	rtext = string(left)
	return
}
