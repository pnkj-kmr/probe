package profile

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"probe/apiserver/handler"
	M "probe/model"
)

func (api *R) saveEnvironmentVariable(env M.Env) (err error) {
	db, ok := api.DB.Get(string(M.ENV))
	if !ok {
		return fmt.Errorf("NO_ENV_VARIABLE")
	}

	var errs []error
	envMap := handler.StructToMap(env, "db")
	// fmt.Println("--->env", env)
	// fmt.Println("----->envMap", envMap)
	for k, v := range envMap {
		dd, err := json.Marshal(v)
		if err == nil {
			err = db.Create(k, dd)
			if err != nil {
				errs = append(errs, err)
				slog.Info("UNABLE TO SAVE ENV VARIABLE")
				continue
			}
		} else {
			errs = append(errs, err)
			continue
		}
	}

	var _err string = ""
	for _, e := range errs {
		_err += e.Error()
	}
	if _err != "" {
		return fmt.Errorf("ENV_ERROR_FOUND: %s", _err)
	}

	return
}

func (api *R) getEnvironmentVariable() (data map[string]string, err error) {
	db, ok := api.DB.Get(string(M.ENV))
	if !ok {
		return data, fmt.Errorf("NO_ENV_VARIABLE")
	}

	r, err := db.FindAll()
	if err != nil {
		return data, err
	}

	data = make(map[string]string)
	for k, v := range r {
		data[k] = string(v)
	}
	return
}
