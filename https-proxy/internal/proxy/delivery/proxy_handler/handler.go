package proxy_handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"

	http_usecase "http-proxy/internal/proxy/usecase/http"
)

type proxyHandler struct {
	httpUsecase http_usecase.HttpUsecase
	logger      *logrus.Logger
}

func NewProxyHandler(logger *logrus.Logger) *proxyHandler {
	return &proxyHandler{
		logger: logger,
	}
}
func (h *proxyHandler) WarnLog(format string, args ...interface{}) {
	h.logger.Warnf(format, args...)
}
func (h *proxyHandler) ErrLog(format string, args ...interface{}) {
	h.logger.Errorf(format, args...)
}
func (h *proxyHandler) InfoLog(format string, args ...interface{}) {
	h.logger.Infof(format, args...)
}
func (h *proxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	h.InfoLog("Request as dump: %v", string(dump))
	w.Write([]byte("Hello, bro"))

	//if r.URL.Scheme != "http" {
	//	h.InfoLog("unsupported protocol scheme = %v\n", r.URL.Scheme)
	//	http.Error(w, fmt.Sprintf("unsupported protocol scheme = %s\n", r.URL.Scheme), http.StatusBadRequest)
	//	return
	//}
	proxyRequest, err := http.NewRequest(r.Method, r.RequestURI, r.Body)
	if h.httpUsecase.Handle(w, r, proxyRequest) != nil {
		h.ErrLog("httpUsecase err %v\n", err)
		return

	}

}
