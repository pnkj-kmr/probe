package profile

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"probe/apiserver/handler"
	M "probe/model"
	"probe/safemap"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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
	api.mux.Get("/agent/info", api.ProbeConfiguration)
	api.mux.Post("/agent/token", api.ProbeConfigurationUpdate)

	api.mux.Post("/text/encrypt", api.TextEncrypt)
	api.mux.Post("/text/decrypt", api.TextDecrypt)

	api.mux.Get("/env", api.GetEnv)
	api.mux.Post("/env", api.SaveEnv)

	api.mux.Get("/process/status", api.GetProcessStatus)

	return api.mux
}

// ProbeConfiguration godoc
// @Summary Get Agent Info
// @Description Probe configuration detail
// @Tags Profile
// @Accept  json
// @Produce  json
// @Success 200 {object}  nil
// @Router /api/profile/agent/info [get]
func (api *R) ProbeConfiguration(w http.ResponseWriter, r *http.Request) {
	config, err := api.getAgentConfiguration()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	data, err := json.Marshal(config)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// ProbeConfigurationUpdate godoc
// @Summary Update Agent Token
// @Description Probe configuration update
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param input body TokenPayload true "Token Payload"
// @Success 201 {object}  TokenPayload
// @Router /api/profile/agent/token [post]
func (api *R) ProbeConfigurationUpdate(w http.ResponseWriter, r *http.Request) {
	token := TokenPayload{}
	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		slog.Info("Token payload issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	err = api.updateAgentToken(token.Token)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	data, err := json.Marshal(token)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	w.Write(data)
}

// TextEncrypt godoc
// @Summary Encryting the text
// @Description encryption
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param input body TextPayload true "Request Payload"
// @Success 201 {object}  TextPayload
// @Router /api/profile/text/encrypt [post]
func (api *R) TextEncrypt(w http.ResponseWriter, r *http.Request) {
	payload := TextPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Info("text payload issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	cryptText, err := handler.Encrypt([]byte(payload.T))
	if err != nil {
		slog.Info("encypting issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	payload2 := TextPayload{T: string(cryptText)}
	data, err := json.Marshal(payload2)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	w.Write(data)
}

// TextDecrypt godoc
// @Summary Decrypt a text
// @Description decryption of payload
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param input body TextPayload true "Request Payload"
// @Success 201 {object}  TextPayload
// @Router /api/profile/text/decrypt [post]
func (api *R) TextDecrypt(w http.ResponseWriter, r *http.Request) {
	payload := TextPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Info("text payload issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	cryptText, err := handler.Decrypt([]byte(payload.T))
	if err != nil {
		slog.Info("decruption issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	payload2 := TextPayload{T: string(cryptText)}
	data, err := json.Marshal(payload2)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	w.Write(data)
}

// GetEnv godoc
// @Summary Get Env Info
// @Description Probe env detail
// @Tags Profile
// @Accept  json
// @Produce  json
// @Success 200 {object}  nil
// @Router /api/profile/env [get]
func (api *R) GetEnv(w http.ResponseWriter, r *http.Request) {
	config, err := api.getEnvironmentVariable()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	data, err := json.Marshal(config)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// SaveEnv godoc
// @Summary Save Env Info
// @Description Probe env detail
// @Tags Profile
// @Accept  json
// @Produce  json
// @Param input body M.Env true "Request Payload"
// @Success 201 {object}  M.Env
// @Router /api/profile/env [post]
func (api *R) SaveEnv(w http.ResponseWriter, r *http.Request) {
	payload := M.Env{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		slog.Info("env payload issue", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	v := validator.New()
	err = v.Struct(payload)
	if err != nil {
		slog.Info("envronment validation", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	err = api.saveEnvironmentVariable(payload)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	data, err := json.Marshal(payload)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}

// GetProcessStatus godoc
// @Summary Get Processes Info
// @Description Probe thread process detail
// @Tags Profile
// @Accept  json
// @Produce  json
// @Success 200 {array}  M.ProcessBeat
// @Router /api/profile/process/status [get]
func (api *R) GetProcessStatus(w http.ResponseWriter, r *http.Request) {
	detail, err := api.getProcessStatus()
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	data, err := json.Marshal(detail)
	if err != nil {
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusOK)
	w.Write(data)
}
