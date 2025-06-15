package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func spaHandler(fsys fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(fsys))

	return func(w http.ResponseWriter, r *http.Request) {
		requestPath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if requestPath == "" {
			requestPath = "index.html"
		}

		// Check if the requested file exists
		f, err := fsys.Open(requestPath)
		if err == nil {
			// File exists, serve it
			f.Close() // We just needed to check existence
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist, serve index.html
		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	}
}
