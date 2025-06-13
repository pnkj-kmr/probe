package credential

import (
	"encoding/json"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// TelnetCredentialHandler godoc
// @Summary      Save Login Credential
// @Description  help to save login profile or device profile creds
// @Tags         Credential
// @Accept       json
// @Produce      json
// @Param input body M.AuthTELNET true "Request payload"
// @Success      201  {object}  M.AuthTELNET
// @Router       /api/credential/telnet [post]
func (api *R) save_telnet(w http.ResponseWriter, r *http.Request) {
	data := &M.AuthTELNET{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	v := validator.New()
	err = v.Struct(data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	d, err := json.Marshal(data)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	err = api.DB.Create(data.LoginProfileId, d)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	w.Write([]byte(d))
	render.Status(r, http.StatusCreated)
}
