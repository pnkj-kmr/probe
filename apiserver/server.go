package api

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"probe/apiserver/handler"
	"probe/apiserver/handler/credential"
	"probe/apiserver/handler/dashboard"
	"probe/apiserver/handler/poll"
	"probe/apiserver/handler/profile"

	_ "probe/docs"
	M "probe/model"
	"probe/safemap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//go:embed static/*
var staticFiles embed.FS

/**
API Server
*/

type Server struct {
	port   int
	sender M.Sender[any]
	db     *safemap.SafeMap[string, M.DB]
}

func New(sender M.Sender[any], db *safemap.SafeMap[string, M.DB]) *Server {
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

	// Add CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // or specific domains
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by major browsers
	}))

	// Serve static files
	uiFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	r.Handle("/*", handler.SPAHandler(uiFS))

	// Swagger UI
	r.Get("/docs/*", httpSwagger.WrapHandler)

	// Group: /api
	r.Route("/api", func(api chi.Router) {
		api.Use(AddDefaultHeader)
		// Group: /
		api.Get("/", server.Ping)

		// Group: /api/poll
		api.Route("/poll", func(api2 chi.Router) {
			api2.Mount("/icmp", poll.NewRouter(string(M.ICMP), server.db).Mux())
			api2.Mount("/snmp", poll.NewRouter(string(M.SNMP), server.db).Mux())
		})

		// Group: /api/credential
		api.Route("/credential", func(api2 chi.Router) {
			credDB, _ := server.db.Get(string(M.CRED))
			api2.Mount("/snmp", credential.NewRouter(M.SNMPType, credDB).Mux())
			api2.Mount("/ssh", credential.NewRouter(M.SSHType, credDB).Mux())
			api2.Mount("/telnet", credential.NewRouter(M.TELNETType, credDB).Mux())
			api2.Mount("/http", credential.NewRouter(M.HTTPType, credDB).Mux())
			api2.Mount("/sftp", credential.NewRouter(M.SFTPType, credDB).Mux())
		})

		// Group: /api/dashboard
		api.Mount("/dashboard", dashboard.NewRouter(server.db).Mux())

		// Group: /api/profile
		api.Mount("/profile", profile.NewRouter(server.db).Mux())

	})

	return r
}

// PingHandler godoc
// @Summary      Ping
// @Description  Responds with pong
// @Tags         Health
// @Produce      json
// @Success      200  {string}  string "pong"
// @Router       /api/ [get]
func (server *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func AddDefaultHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
