package uiserver

import (
	"github.com/go-chi/chi"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UiServer(r chi.Router, path string) {
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "ui")
	root := http.Dir(filesDir)
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
