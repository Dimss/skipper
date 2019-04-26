package sunkey

type Node struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

type SunkeyData struct {
	Nodes []Node `json:"nodes"`
	Links []Link `json:"links"`
}

