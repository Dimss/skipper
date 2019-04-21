package rbac

import (
	"encoding/json"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := GetUsers()
	responsePayload := ResponsePayload{Status: "ok", Data: users}
	if response, err := json.Marshal(responsePayload); err == nil {
		w.Write(response)
	}
}

func GetRolesHandler(w http.ResponseWriter, r *http.Request) {
	//nodes := []Nodes{
	//	{"Get", "get"},
	//	{"Watch", "watch"},
	//	{"List", "list"},
	//	{"Default", "default"},
	//	{"Nodes", "nodes"},
	//	{"Services", "services"},
	//	{"Endpoints", "endpoints"},
	//	{"Pods", "pods"},
	//	{"myRole", "MyRole"},
	//}
	//links := []Links{
	//	{0, 3, 1},
	//	{1, 3, 1},
	//	{2, 3, 1},
	//	{3, 4, 1},
	//	{3, 5, 1},
	//	{3, 6, 1},
	//	{0, 8, 1},
	//	{1, 8, 1},
	//	{2, 8, 1},
	//	{8, 7, 1},
	//	{8, 4, 1},
	//	{8, 5, 1},
	//	{8, 6, 1},
	//	{8, 7, 1},
	//}
	//
	//
	//data := SunkeyData{nodes, links}
	roles := GetRoles()
	responsePayload := ResponsePayload{Status: "ok", Data: roles}
	if response, err := json.Marshal(responsePayload); err == nil {
		w.Write(response)
	}
}

func GetNamespacesHandler(w http.ResponseWriter, r *http.Request) {
	users := GetNamespaces()
	responsePayload := ResponsePayload{Status: "ok", Data: users}
	if response, err := json.Marshal(responsePayload); err == nil {
		w.Write(response)
	}
}
