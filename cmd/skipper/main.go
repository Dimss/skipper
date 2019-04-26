package main

import (
	"github.com/dimss/skipper/pkg/namespaces"
	"github.com/dimss/skipper/pkg/roles"
	"github.com/dimss/skipper/pkg/uiserver"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Get("/roles", roles.GetRolesHandler)
	r.Get("/namespaces", namespaces.GetNamespacesHandler)
	uiserver.UiServer(r, "/ui")
	http.ListenAndServe(":3000", r)
}
