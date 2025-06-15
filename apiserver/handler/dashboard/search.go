package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"
	"strings"

	"github.com/go-chi/render"
)

// GetModelList godoc
// @Summary Get search model name
// @Description stats db models
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Success 200 {object} ModelList
// @Router /api/dashboard/model [get]
func (api *R) GetModelList(w http.ResponseWriter, r *http.Request) {
	var list ModelList

	list = append(list, SearchModel{L: strings.ToUpper(fmt.Sprintf("%s", M.CRED)), V: fmt.Sprintf("%s", M.CRED)})
	// icmp
	list = append(list, SearchModel{L: strings.ToUpper(fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_60)), V: fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_60)})
	list = append(list, SearchModel{L: strings.ToUpper(fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_300)), V: fmt.Sprintf("%s/%d", M.ICMP, M.INTERVAL_300)})
	// snmp
	list = append(list, SearchModel{L: strings.ToUpper(fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_60)), V: fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_60)})
	list = append(list, SearchModel{L: strings.ToUpper(fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_300)), V: fmt.Sprintf("%s/%d", M.SNMP, M.INTERVAL_300)})

	data, err := json.Marshal(list)
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
	}

	data, err := json.Marshal(stat)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}
