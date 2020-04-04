package registry

// 服务
type Service struct {
	Name  string  `json:"name"`
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	Id     int    `json:"id"`
	Ip     string `json:"ip"`
	Port   int    `json:"port"`
	Weigth int    `json:"weigth"`
}
