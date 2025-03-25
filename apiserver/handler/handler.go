package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	M "probe/model"

	"github.com/go-chi/render"
)

type handle struct {
	sender M.Sender[[]byte]
	// receiver M.Receiver[[]byte]
}

func New(s M.Sender[[]byte]) *handle {
	return &handle{
		sender: s,
		// receiver: r,
	}
}

func (h *handle) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func (h *handle) PaylodICMP(w http.ResponseWriter, r *http.Request) {

	data := &M.InICMP{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// article := data.Article
	// dbNewArticle(article)
	fmt.Println("data ---> ", data)
	d, _ := json.Marshal(data)
	h.sender.Send(d)
	w.Write(d)
	render.Status(r, http.StatusCreated)

}
