package profile

import (
	"fmt"
	M "probe/model"
	"time"
)

var configurationKeys = []string{
	"token",
	"id",
	"sid",
	"name",
	"ip_address",
	"hostname",
	"osname",
	"role",
	"secret_key",
	"org_id",
	"created",
	"modified",
}

type TokenPayload struct {
	Token string `json:"token"`
}

type TextPayload struct {
	T string `json:"text"`
}

func (api *R) getAgentConfiguration() (data map[string]any, err error) {
	db, ok := api.DB.Get(string(M.CONFIG))
	if !ok {
		return nil, fmt.Errorf("NO_CONFIGURATION")
	}

	data = make(map[string]any)
	for _, k := range configurationKeys {
		r, _ := db.Find(k)
		data[k] = string(r)
	}

	return
}

func (api *R) updateAgentToken(token string) (err error) {
	var key = "token"
	db, ok := api.DB.Get(string(M.CONFIG))
	if !ok {
		return fmt.Errorf("NO_CONFIGURATION")
	}

	// checking token exists for not
	r, err := db.Find(key)
	if err == nil {
		if len(r) > 1 {
			return fmt.Errorf("AGENT_TOKEN_EXISTS")
		}
	}

	// crating new token into system
	err = db.Create(key, []byte(token))
	if err != nil {
		return err
	}
	_now := time.Now().UTC().String()
	db.Create("created", []byte(_now))
	db.Create("modified", []byte(_now))

	// TODO
	// need to create more keys into db

	return
}
