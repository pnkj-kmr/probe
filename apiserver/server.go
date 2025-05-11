package api

import (
	"fmt"
	"net/http"
	"probe/apiserver/handler/poll"
	M "probe/model"
	"probe/safemap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/**
API Server
*/

type Server struct {
	port   int
	sender M.Sender[any]
	db     *safemap.SafeMap[int, M.DB]
}

func New(sender M.Sender[any], db *safemap.SafeMap[int, M.DB]) *Server {
	return &Server{
		port:   3000,
		sender: sender,
		db:     db,
	}
}

// Run helps to spin the API server
func (server *Server) Run() error {
	r := server.newRouter()
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", server.port), r)
}

func (server *Server) newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Group: /
	r.Get("/", server.Ping)

	// Group: /api
	r.Route("/api", func(api chi.Router) {

		// Group: /api/icmp
		api.Mount("/icmp", poll.NewRouter(M.ICMP, server.db).Mux())

		// // Group: /api/snmp
		api.Mount("/snmp", poll.NewRouter(M.SNMP, server.db).Mux())

	})

	return r
}

func (server *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
