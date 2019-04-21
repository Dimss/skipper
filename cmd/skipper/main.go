package main

import (
	"github.com/dimss/skipper/pkg/rbac"
	"github.com/dimss/skipper/pkg/uiserver"
	"github.com/go-chi/chi"
	"net/http"
)

//https://bl.ocks.org/wvengen/cab9b01816490edb7083
func main() {

	r := chi.NewRouter()
	r.Get("/users", rbac.GetUsersHandler)
	r.Get("/roles", rbac.GetRolesHandler)
	r.Get("/namespaces", rbac.GetNamespacesHandler)
	uiserver.UiServer(r, "/ui")
	http.ListenAndServe(":3000", r)
}
