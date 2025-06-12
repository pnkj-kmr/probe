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

type R struct {
	mux *chi.Mux

	ID string
	DB *safemap.SafeMap[string, M.DB]
}

func NewRouter(id string, db *safemap.SafeMap[string, M.DB]) *R {
	return &R{
		mux: chi.NewRouter(),
		ID:  id,
		DB:  db,
	}
}

func (api *R) Mux() *chi.Mux {
	api.mux.Post("/", api.Save)
	api.mux.Get("/{pollPeriod}/{id}", api.Get)
	api.mux.Delete("/{pollPeriod}/{id}", api.Delete)
	api.mux.Get("/{pollPeriod}", api.Count)
	api.mux.Delete("/{pollPeriod}", api.DeleteAll)

	return api.mux
}

func (api *R) Save(w http.ResponseWriter, r *http.Request) {
	switch M.ProbeType(api.ID) {
	case M.ICMP:
		api.save_icmp(w, r)
	case M.SNMP:
		api.save_snmp(w, r)
	default:
		render.Render(w, r, handler.ErrInvalidRequest(handler.ErrNoData))
	}
}

// DeletePollConfig godoc
// @Summary Poll Config Delete
// @Description delete poll config
// @Tags Poll
// @Accept  json
// @Produce  json
// @Param type path string true "type: icmp / snmp"
// @Param pollPeriod path string true "Poll Period - like: 60/300"
// @Param id path string true "CI ID of resource"
// @Success 200 {string} string "record"
// @Router /api/poll/{type}/{pollPeriod}/{id} [delete]
func (api *R) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.GetDB(pollPeriod)
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

// DeleteAllPollCount godoc
// @Summary Delete All Poll Config
// @Description deleting the poll config records
// @Tags Poll
// @Accept  json
// @Produce  json
// @Param type path string true "type: icmp / snmp"
// @Param pollPeriod path string true "Poll Period - like: 60/300"
// @Success 200 {string} string "record"
// @Router /api/poll/{type}/{pollPeriod} [delete]
func (api *R) DeleteAll(w http.ResponseWriter, r *http.Request) {
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.GetDB(pollPeriod)
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

// GetTotalPollCount godoc
// @Summary Poll Config Total Count
// @Description get totel record of poll config
// @Tags Poll
// @Accept  json
// @Produce  json
// @Param type path string true "type: icmp / snmp"
// @Param pollPeriod path string true "Poll Period - like: 60/300"
// @Success 200 {string} string "record"
// @Router /api/poll/{type}/{pollPeriod} [get]
func (api *R) Count(w http.ResponseWriter, r *http.Request) {
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.GetDB(pollPeriod)
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
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// GetPollConfig godoc
// @Summary Poll Config Get View
// @Description get poll config data
// @Tags Poll
// @Accept  json
// @Produce  json
// @Param type path string true "type: icmp / snmp"
// @Param pollPeriod path string true "Poll Period - like: 60/300"
// @Param id path string true "CI ID of resource"
// @Success 200 {string} string "poll config"
// @Router /api/poll/{type}/{pollPeriod}/{id} [get]
func (api *R) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	pollPeriod := chi.URLParam(r, "pollPeriod")
	db, err := api.GetDB(pollPeriod)
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

func (api *R) GetDB(pollPeriod any) (M.DB, error) {
	var dbId string
	switch x := pollPeriod.(type) {
	case int:
		dbId = api.ID + "/" + strconv.Itoa(x)
	case string:
		dbId = api.ID + "/" + x
	}
	// fmt.Println("dbId --->", dbId)
	db, ok := api.DB.Get(dbId)
	if !ok {
		return nil, handler.ErrDB
	}
	return db, nil
}
