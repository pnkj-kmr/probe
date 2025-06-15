package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"
	"probe/safemap"
	"strings"

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
	api.mux.Get("/data", api.GetStats)
	api.mux.Get("/model", api.GetModelList)

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

	// credential profile count
	t := api.resCount(string(M.CRED))
	stats = append(stats, PollStat{Name: strings.ToUpper(string(M.CRED)), Total: t})

	// getting icmp stats
	t = api.resourceCount(M.ICMP)
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
	var stats []PollStat
	var t, t2 uint64

	// icmp stats
	t = api.resCount(fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_60))
	t2 = api.resCount(fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_60))
	stats = append(stats, PollStat{Name: "ICMP_60", Total: t, Polled: t2})
	t = api.resCount(fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_300))
	t2 = api.resCount(fmt.Sprintf("%s/%d_stat", M.ICMP, M.INTERVAL_300))
	stats = append(stats, PollStat{Name: "ICMP_300", Total: t, Polled: t2})

	// snmp stats
	t = api.resCount(fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_60))
	t2 = api.resCount(fmt.Sprintf("%s/%d_stat", M.SNMP, M.INTERVAL_60))
	stats = append(stats, PollStat{Name: "SNMP_60", Total: t, Polled: t2})
	t = api.resCount(fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_300))
	t2 = api.resCount(fmt.Sprintf("%s/%d_stat", M.SNMP, M.INTERVAL_300))
	stats = append(stats, PollStat{Name: "SNMP_300", Total: t, Polled: t2})

	data, err := json.Marshal(stats)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}
