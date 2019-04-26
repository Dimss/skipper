package main

import (
	"github.com/dimss/skipper/pkg/rbac"
	"github.com/dimss/skipper/pkg/uiserver"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
)

//https://bl.ocks.org/wvengen/cab9b01816490edb7083
func main() {

	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Get("/users", rbac.GetUsersHandler)
	r.Get("/roles", rbac.GetRolesHandler)
	r.Get("/namespaces", rbac.GetNamespacesHandler)

	uiserver.UiServer(r, "/ui")
	http.ListenAndServe(":3000", r)
}
