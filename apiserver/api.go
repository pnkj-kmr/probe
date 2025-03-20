package api

/**
API Server
*/

type API struct {
}

func New() *API {
	return &API{}
}

// Run helps to spin the API server
func (api *API) Run() error {
	return nil
}
