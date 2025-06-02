package api

import (
	"fmt"
	"net/http"
	"probe/apiserver/handler/auth"
	"probe/apiserver/handler/poll"

	_ "probe/docs"
	M "probe/model"
	"probe/safemap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Swagger UI
	r.Get("/docs/*", httpSwagger.WrapHandler)
	// r.Get("/docs/*", httpSwagger.WrapHandler(liteFiles.Handler))
	// // Serve your swagger.json file (optional)
	// r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./docs/swagger.json")
	// })

	// Group: /api
	r.Route("/api", func(api chi.Router) {

		// Group: /api/poll
		api.Route("/poll", func(api2 chi.Router) {
			api2.Mount("/icmp", poll.NewRouter(M.ICMP, server.db).Mux())
			api2.Mount("/snmp", poll.NewRouter(M.SNMP, server.db).Mux())
		})

		// Group: /api/profile
		api.Route("/profile", func(api2 chi.Router) {
			authDB, _ := server.db.Get(M.AUTH_PROFILE)
			api2.Mount("/snmp", auth.NewRouter(M.SNMPType, authDB).Mux())
			api2.Mount("/ssh", auth.NewRouter(M.SSHType, authDB).Mux())
			api2.Mount("/telnet", auth.NewRouter(M.TELNETType, authDB).Mux())
			api2.Mount("/http", auth.NewRouter(M.HTTPType, authDB).Mux())
			api2.Mount("/sftp", auth.NewRouter(M.SFTPType, authDB).Mux())
		})

	})

	return r
}

// PingHandler godoc
// @Summary      Ping
// @Description  Responds with pong
// @Tags         Health
// @Produce      json
// @Success      200  {string}  string "pong"
// @Router       / [get]
func (server *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
