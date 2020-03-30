package registry

// 服务
type Service struct {
	Name string
	Node []*Node
}

type Node struct {
	Id   int
	Ip   string
	Port int
}
