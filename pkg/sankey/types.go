package sankey

type Sankey interface {
	CreateGraphData(k8sObjects interface{}) *GraphData
}

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

type GraphData struct {
	Nodes         []Node `json:"nodes"`
	Links         []Link `json:"links"`
	nodesIndexMap map[string]int
	k8sObjects    interface{}
}
