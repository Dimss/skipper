package uiserver

import (
	"github.com/go-chi/chi"
	"github.com/gobuffalo/packr/v2"
	"net/http"
	"strings"
)

func UiServer(r chi.Router, path string) {

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	box := packr.New("someBoxName", "./../../ui")
	fs := http.StripPrefix(path, http.FileServer(box))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/ui/" || !box.Has(r.URL.Path) {
			r.RequestURI = "/ui/app-index.html"
			r.URL.Path = "/ui/app-index.html"
		}
		fs.ServeHTTP(w, r)
	}))
}
