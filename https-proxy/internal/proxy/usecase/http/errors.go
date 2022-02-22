package http_usecase

import "errors"

var (
	NilError    = errors.New("nil argument in func params")
	NilResponse = errors.New("nil response")
)
