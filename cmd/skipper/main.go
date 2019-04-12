package main

import (
	"github.com/dimss/skipper/pkg/rbac"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Get("/users", rbac.GetUsersHandler)
	r.Get("/roles", rbac.GetRolesHandler)
	r.Get("/namespaces", rbac.GetNamespacesHandler)
	http.ListenAndServe(":3000", r)
}
