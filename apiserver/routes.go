package api

import (
	"probe/apiserver/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Server) apiRouter() *chi.Mux {
	r := chi.NewRouter()

	h := handler.New(api)

	r.Get("/", h.Ping)
	r.Post("/icmp", h.PaylodICMP)
	r.Get("/icmp", h.PaylodICMP)
	return r
}

func (api *Server) newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api", api.apiRouter())

	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	return r
}
