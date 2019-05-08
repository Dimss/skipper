package main

import (
	"github.com/dimss/skipper/pkg/clusterrole"
	"github.com/dimss/skipper/pkg/clusterrolebinding"
	"github.com/dimss/skipper/pkg/namespaces"
	"github.com/dimss/skipper/pkg/rolebinding"
	"github.com/dimss/skipper/pkg/roles"
	"github.com/dimss/skipper/pkg/uiserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Get("/roles", roles.GetRolesHandler)
	r.Get("/bindings", rolebinding.GetRolesBindingsHandler)
	r.Get("/namespaces", namespaces.GetNamespacesHandler)
	r.Get("/clusterroles", clusterrole.GetClusterRolesHandler)
	r.Get("/clusterrolesbindings", clusterrolebinding.GetClusterRolesBindingsHandler)
	uiserver.UiServer(r, "/ui")
	http.ListenAndServe(":3001", r)
}
