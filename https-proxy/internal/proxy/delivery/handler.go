package delivery

import (
	"net/http"
)

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	WarnLog(msg string)
	ErrLog(msg string)
	InfoLog(format string, args ...interface{})
}
