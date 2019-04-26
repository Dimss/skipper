package roles

import (
	"encoding/json"
	"net/http"
)

func GetRolesHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	sunkeyData := GetRoles(ns)
	if response, err := json.Marshal(sunkeyData); err == nil {
		w.Write(response)
	}
}
