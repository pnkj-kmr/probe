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
		dd, err := json.Marshal(M.Record{K: k, V: v})
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

func (api *R) getEnvironmentVariable() (data []M.Record, err error) {
	db, ok := api.DB.Get(string(M.ENV))
	if !ok {
		return data, fmt.Errorf("NO_ENV_VARIABLE")
	}

	for v := range db.Receive() {
		var record M.Record
		err := json.Unmarshal(v, &record)
		if err != nil {
			return data, err
		}
		data = append(data, record)
	}
	return
}
