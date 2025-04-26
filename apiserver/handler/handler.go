package handler

import (
	"net/http"
	M "probe/model"

	"github.com/go-chi/render"
)

type handle struct {
	sender M.Sender[any]
	// receiver M.Receiver[[]byte]
}

func New(s M.Sender[any]) *handle {
	return &handle{
		sender: s,
		// receiver: r,
	}
}

func (h *handle) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (h *handle) PaylodICMP(w http.ResponseWriter, r *http.Request) {

	data := &M.ICMPReq{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	h.sender.Send(data)
	w.Write([]byte("request has been sent to queue."))
	render.Status(r, http.StatusCreated)

}

func (h *handle) GetICMP(w http.ResponseWriter, r *http.Request) {

	data := &M.ICMPReq{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	h.sender.Send(data)
	w.Write([]byte("request has been sent to queue."))
	render.Status(r, http.StatusCreated)

}
