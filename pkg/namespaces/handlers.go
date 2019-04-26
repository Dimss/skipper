package namespaces

import (
	"encoding/json"
	"net/http"
)



func GetNamespacesHandler(w http.ResponseWriter, r *http.Request) {
	users := GetNamespaces()
	if response, err := json.Marshal(users); err == nil {
		w.Write(response)
	}
}
