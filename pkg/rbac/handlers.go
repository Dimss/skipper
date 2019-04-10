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
	users := GetRoles()
	responsePayload := ResponsePayload{Status: "ok", Data: users}
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