package poll

import (
	M "probe/model"

	"github.com/go-playground/validator/v10"
)

func (api *R) validateAuthProfile(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	db, ok := api.DB.Get(string(M.CRED))
	// fmt.Println("---->", db)
	if !ok {
		return ok
	}
	data, err := db.Find(value)
	// fmt.Println("---->data", data)
	if err != nil {
		return false
	}
	if len(data) == 0 {
		return false
	}
	return true
}
