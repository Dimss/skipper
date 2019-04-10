package rbac

type ResponsePayload struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Users struct {
	Users []string `json:"users"`
}
