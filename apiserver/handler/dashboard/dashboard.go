package dashboard

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"
	"probe/safemap"

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
	api.mux.Get("/system", api.systemStats)
	api.mux.Get("/resource", api.resourceStats)
	// api.mux.Delete("/{pollPeriod}/{id}", api.delete)
	// api.mux.Get("/{pollPeriod}", api.count)
	// api.mux.Delete("/{pollPeriod}", api.deleteAll)

	return api.mux
}

// SystemStats godoc
// @Summary Get system stats
// @Description cpu, memory and disk usage
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {string} string "system stats"
// @Router /api/dashboard/system [get]
func (api *_router) systemStats(w http.ResponseWriter, r *http.Request) {
	stats := getSystemStats()
	data, err := json.Marshal(stats)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// ResourceStats godoc
// @Summary Get resource stats
// @Description total resource which aligned to poll
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {string} string "system stats"
// @Router /api/dashboard/resource [get]
func (api *_router) resourceStats(w http.ResponseWriter, r *http.Request) {
	var data []byte

	// TODO
	// need to group all poll type stats
	//

	render.Status(r, http.StatusOK)
	w.Write(data)
}
