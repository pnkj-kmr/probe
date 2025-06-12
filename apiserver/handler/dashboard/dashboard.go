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

type R struct {
	mux *chi.Mux

	DB *safemap.SafeMap[string, M.DB]
}

func NewRouter(db *safemap.SafeMap[string, M.DB]) *R {
	return &R{
		mux: chi.NewRouter(),
		DB:  db,
	}
}

func (api *R) Mux() *chi.Mux {
	api.mux.Get("/system", api.SystemStats)
	api.mux.Get("/resource", api.ResourceStats)
	api.mux.Get("/poll", api.PollStats)

	return api.mux
}

// SystemStats godoc
// @Summary Get system stats
// @Description cpu, memory and disk usage
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object}  SystemStats
// @Router /api/dashboard/system [get]
func (api *R) SystemStats(w http.ResponseWriter, r *http.Request) {
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
// @Success 200 {object}  ResourceStats
// @Router /api/dashboard/resource [get]
func (api *R) ResourceStats(w http.ResponseWriter, r *http.Request) {
	var stats []PollStat
	var total uint64

	// getting icmp stats
	t := api.resourceCount(M.ICMP)
	stats = append(stats, PollStat{Name: "ICMP", Total: t})
	total = total + t

	t = api.resourceCount(M.SNMP)
	stats = append(stats, PollStat{Name: "SNMP", Total: t})
	total = total + t

	// TODO
	// add more here

	resourceStat := ResourceStats{
		Stats: stats,
		Total: total,
	}
	data, err := json.Marshal(resourceStat)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// PollStats godoc
// @Summary Get polling stats
// @Description helps to get the current polling status along with poll result
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object}  PollingStats
// @Router /api/dashboard/poll [get]
func (api *R) PollStats(w http.ResponseWriter, r *http.Request) {

}
