package api

import (
	"fmt"
	"net/http"
	M "probe/model"
)

/**
API Server
*/

type Server struct {
	port   int
	sender M.Sender[[]byte]
	// receiver M.Receiver[[]byte]
}

func New(sender M.Sender[[]byte]) *Server {
	return &Server{
		port:   3000,
		sender: sender,
		// receiver: receiver,
	}
}

// Run helps to spin the API server
func (api *Server) Run() error {
	r := api.newRouter()
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", api.port), r)
}

func (api *Server) Send(msg []byte) error {
	return api.sender.Send(msg)
}

// func (api *Server) Receive() []byte {
// 	return api.receiver.Receive()
// }
