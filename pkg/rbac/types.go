package rbac

type ResponsePayload struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Users struct {
	Users []string `json:"users"`
}

type Nodes struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Links struct {
	Source int `json:"source"`
	Target int `json:"target"`
	Value  int `json:"value"`
}

type SunkeyData struct {
	Nodes []Nodes `json:"nodes"`
	Links []Links `json:"links"`
}
