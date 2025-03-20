package probe

import "errors"

var (
	ErrDB     = errors.New("NO DB CONFIGURATION")
	ErrPOLLER = errors.New("NO POLLER CONFIGURATION")
)
