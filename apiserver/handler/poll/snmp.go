package poll

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// PollHandler godoc
//
// @Summary      Polling Status update
// @Description  APIs helps to write poll parameter into db
// @Security     BearerAuth
// @Tags         Poll
// @Accept       json
// @Produce      json
// @Success      200  {object}  M.SNMPReq
// @Router       /api/poll/snmp [post]
func (api *R) save_snmp(w http.ResponseWriter, r *http.Request) {
	data := &M.SNMPReq{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	v := validator.New()
	v.RegisterValidation("isValidAuth", api.validateAuthProfile)
	err = v.Struct(data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	db, err := api.GetDB(data.Params.PollPeriod)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	d, err := json.Marshal(data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	err = db.Create(data.Cid, d)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	w.Write([]byte(d))
}
