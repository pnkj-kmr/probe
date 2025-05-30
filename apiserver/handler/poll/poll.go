package poll

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"
	"probe/safemap"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type _router struct {
	id  int
	mux *chi.Mux
	db  *safemap.SafeMap[int, M.DB]
}

func NewRouter(id int, db *safemap.SafeMap[int, M.DB]) *_router {
	return &_router{
		id:  id,
		mux: chi.NewRouter(),
		db:  db,
	}
}

func (api *_router) Mux() *chi.Mux {
	api.mux.Post("/", api.save)
	api.mux.Get("/{pollPeriod}/{id}", api.get)
	api.mux.Delete("/{pollPeriod}/{id}", api.delete)
	api.mux.Get("/{pollPeriod}", api.count)
	api.mux.Delete("/{pollPeriod}", api.deleteAll)

	return api.mux
}

func (api *_router) save(w http.ResponseWriter, r *http.Request) {
	switch api.id {
	case M.ICMP:
		api.save_icmp(w, r)
	case M.SNMP:
		api.save_snmp(w, r)
	default:
		render.Render(w, r, handler.ErrInvalidRequest(handler.ErrNoData))
	}
}

func (api *_router) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.getDB(pollPeriod)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	err = db.Delete(id)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (api *_router) deleteAll(w http.ResponseWriter, r *http.Request) {
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.getDB(pollPeriod)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	err = db.DeleteAll()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

func (api *_router) count(w http.ResponseWriter, r *http.Request) {
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.getDB(pollPeriod)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	count := db.Count()
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
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.getDB(pollPeriod)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	data, err := db.Find(id)
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

func (api *_router) getDB(pollPeriod any) (M.DB, error) {
	var dbId int
	switch x := pollPeriod.(type) {
	case int:
		dbId = api.id + x
	case string:
		p, err := strconv.Atoi(x)
		if err != nil {
			p = M.INTERVAL_300
		}
		dbId = api.id + p
	}
	// fmt.Println("dbId --->", dbId)
	db, ok := api.db.Get(dbId)
	if !ok {
		return nil, handler.ErrDB
	}
	return db, nil
}
