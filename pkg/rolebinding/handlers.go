package rolebinding

import (
	"encoding/json"
	"net/http"
)

func GetRolesBindingsHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	sunkeyData := GetBindings(ns)
	if response, err := json.Marshal(sunkeyData); err == nil {
		w.Write(response)
	}
}
