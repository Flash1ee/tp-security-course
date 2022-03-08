package proxy_handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/usecase"
	http_usecase "http-proxy/internal/proxy/usecase/http"
	https_usecase "http-proxy/internal/proxy/usecase/https"
	"http-proxy/pkg/utils"
)

type proxyHandler struct {
	uc     usecase.Usecase
	logger *logrus.Logger
	config *cfg.Config
	conn   *utils.PostgresConn
}

func NewProxyHandler(logger *logrus.Logger, conn *utils.PostgresConn, config *cfg.Config) *proxyHandler {
	return &proxyHandler{
		logger: logger,
		conn:   conn,
		config: config,
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
	var err error
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	h.InfoLog("Request as dump: %v", string(dump))

	if r.Method == http.MethodConnect {
		h.uc, err = https_usecase.NewHttpsUsecase(w, r, h.config, h.logger, h.conn)
	} else {
		h.uc, err = http_usecase.NewHttpUsecase(w, r, h.config, h.logger, h.conn)
	}
	if err != nil {
		h.logger.Error(err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if h.uc.Handle() != nil {
		h.ErrLog("httpUsecase err %v\n", err)
	}
	h.uc.Close()

}
