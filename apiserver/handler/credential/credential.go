package credential

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type R struct {
	mux *chi.Mux

	PType M.ProtocolType
	DB    M.DB
}

func NewRouter(pType M.ProtocolType, db M.DB) *R {
	return &R{
		PType: pType,
		mux:   chi.NewRouter(),
		DB:    db,
	}
}

func (api *R) Mux() *chi.Mux {
	// related routers
	api.mux.Post("/", api.save)
	api.mux.Get("/", api.count)
	api.mux.Delete("/", api.deleteAll)
	api.mux.Get("/{id}", api.get)
	api.mux.Delete("/{id}", api.delete)

	return api.mux
}

func (api *R) save(w http.ResponseWriter, r *http.Request) {
	switch api.PType {
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

// DeleteCredentialByID godoc
// @Summary Delete credential for profile id
// @Description delete one
// @Tags Credential
// @Accept  json
// @Produce  json
// @Param credential_type path string true "type: snmp / ssh / telnet / sftp / http"
// @Param id path string true "credential id (login/device profile id)"
// @Success 200 {string} string "record"
// @Router /api/credential/{credential_type}/{id} [delete]
func (api *R) delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := api.DB.Delete(id)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

// DeleteCredentials godoc
// @Summary Delete all credentials for profile
// @Description delete all
// @Tags Credential
// @Accept  json
// @Produce  json
// @Param credential_type path string true "type: snmp / ssh / telnet / sftp / http"
// @Success 200 {string} string "record"
// @Router /api/credential/{credential_type}/ [delete]
func (api *R) deleteAll(w http.ResponseWriter, r *http.Request) {
	err := api.DB.DeleteAll()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusNoContent)
}

// GetCredentialTotalCount godoc
// @Summary Total records of credential
// @Description total count
// @Tags Credential
// @Accept  json
// @Produce  json
// @Param credential_type path string true "type: snmp / ssh / telnet / sftp / http"
// @Success 200 {string} string "record"
// @Router /api/credential/{credential_type}/ [get]
func (api *R) count(w http.ResponseWriter, r *http.Request) {
	count := api.DB.Count()
	// data, err := json.Marshal(struct{ count uint64 }{count: count})
	data, err := json.Marshal(count)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	w.Write(data)
	render.Status(r, http.StatusOK)
}

// GetCredentialByID godoc
// @Summary Get credential record
// @Description creential profile view
// @Tags Credential
// @Accept  json
// @Produce  json
// @Param credential_type path string true "type: snmp / ssh / telnet / sftp / http"
// @Param id path string true "credential id (login/device profile id)"
// @Success 200 {string} string "record"
// @Router /api/credential/{credential_type}/{id} [get]
func (api *R) get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := api.DB.Find(id)
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
