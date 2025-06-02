package credential

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type _router struct {
	authType M.ProtocolType
	mux      *chi.Mux
	db       M.DB
}

func NewRouter(authType M.ProtocolType, db M.DB) *_router {
	return &_router{
		authType: authType,
		mux:      chi.NewRouter(),
		db:       db,
	}
}

func (api *_router) Mux() *chi.Mux {
	// related routers
	api.mux.Post("/", api.save)
	api.mux.Get("/count", api.count)
	api.mux.Get("/{id}", api.get)
	api.mux.Delete("/{id}", api.delete)
	api.mux.Delete("/", api.deleteAll)

	return api.mux
}

func (api *_router) save(w http.ResponseWriter, r *http.Request) {
	switch api.authType {
	case M.SNMPType:
		api.save_snmp(w, r)
	case M.HTTPType, M.HTTPSType:
		api.save_http(w, r)
	case M.SSHType:
		api.save_ssh(w, r)
	case M.TELNETType:
		api.save_telnet(w, r)
	case M.SFTPType:
		api.save_sftp(w, r)
	default:
		render.Render(w, r, handler.ErrInvalidRequest(handler.ErrNoData))
	}
}

func (api *_router) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := api.db.Delete(id)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (api *_router) deleteAll(w http.ResponseWriter, r *http.Request) {
	err := api.db.DeleteAll()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (api *_router) count(w http.ResponseWriter, r *http.Request) {
	count := api.db.Count()
	// data, err := json.Marshal(struct{ count uint64 }{count: count})
	data, err := json.Marshal(count)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	w.Write(data)
	render.Status(r, http.StatusOK)
}

func (api *_router) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := api.db.Find(id)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	// fmt.Println("----data", data)
	if len(data) == 0 {
		render.Render(w, r, handler.ErrInvalidRequest(handler.ErrNoData))
		return
	}
	w.Write(data)
	render.Status(r, http.StatusOK)
}
