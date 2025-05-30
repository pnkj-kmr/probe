package poll

import (
	"encoding/json"
	"log/slog"
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
// @Tags         Poll
// @Accept       json
// @Produce      json
// @Param input body M.ICMPReq true "ICMP Request body"
// @Success      201  {object}  M.ICMPReq
// @Router       /api/poll/icmp [post]
func (api *_router) save_icmp(w http.ResponseWriter, r *http.Request) {
	data := &M.ICMPReq{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Info("ICMP request decode error", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}

	v := validator.New()
	err = v.Struct(data)
	if err != nil {
		slog.Info("ICMP request submit error", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	db, err := api.getDB(data.Params.PollPeriod)
	if err != nil {
		slog.Info("ICMP request db error", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	d, err := json.Marshal(data)
	if err != nil {
		slog.Info("ICMP request marshal error", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	err = db.Create(data.Cid, d)
	if err != nil {
		slog.Info("ICMP request create error", "err", err)
		render.Render(w, r, handler.ErrInvalidRequest(err))
		return
	}
	render.Status(r, http.StatusCreated)
	w.Write([]byte(d))
}
