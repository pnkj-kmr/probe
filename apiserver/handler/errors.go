package handler

import "errors"

var (
	ErrDB     = errors.New("NO DB CONFIGURATION FOUND")
	ErrNoAuth = errors.New("NO AUTH CONFIGURATION FOUND")
	ErrNoData = errors.New("NO DATA FOUND")
)
