package dashboard

import (
	"encoding/json"
	"fmt"
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
	api.mux.Get("/data", api.GetStats)

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

// GetStats godoc
// @Summary Get saved statistics params
// @Description helps to get the current save stats into system
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param table query string icmp "ICMP collection table"
// @Param cid query string test "CI_ID to query"
// @Success 200 {object} nil
// @Router /api/dashboard/data [get]
func (api *R) GetStats(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	cid := queryParams.Get("cid")
	dbName := queryParams.Get("table")

	var stat = make(map[string]any)
	db, ok := api.DB.Get(dbName)
	if ok {
		r, err := db.Find(cid)
		if err != nil {
			stat["config"] = err.Error()
		} else {
			stat["config"] = api.getData(dbName, r)
		}
	} else {
		render.Render(w, r, handler.ErrInvalidRequest(fmt.Errorf("no db name: %s", dbName)))
		return
	}
	dbName = fmt.Sprintf("%s_stat", dbName)
	db, ok = api.DB.Get(dbName)
	if ok {
		r, err := db.Find(cid)
		if err != nil {
			stat["polled"] = err.Error()
		} else {
			stat["polled"] = api.getData(dbName, r)
		}
	} else {
		render.Render(w, r, handler.ErrInvalidRequest(fmt.Errorf("no db name: %s_stat", dbName)))
		return
	}

	data, err := json.Marshal(stat)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}
