package handler

import (
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func SPAHandler(fsys fs.FS) http.HandlerFunc {
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

var jwtSecret = []byte("your-secret-key")

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims or user context if needed
		next.ServeHTTP(w, r)
	})
}
