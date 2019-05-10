package sankey

type Sankey interface {
	LoadK8SObjects()
	GetK8SObjects() (interface{})
	CreateGraphData()
}

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Value  int    `json:"value"`
}
