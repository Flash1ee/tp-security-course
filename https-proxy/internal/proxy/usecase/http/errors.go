package http_usecase

import "errors"

var (
	NilError    = errors.New("nil argument in func params")
	NilResponse = errors.New("nil response")
	InvalidArg  = errors.New("invalid arg")
	LogicError  = errors.New("logic error")
	DumpError   = errors.New("can not get dump")
)
