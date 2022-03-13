package proxy_history

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/usecase"
	"http-proxy/pkg/utils"
)

type proxyHistory struct {
	uc     usecase.HistoryUsecase
	logger *logrus.Logger
	config *cfg.Config
	conn   *utils.PostgresConn
}

func NewProxyHistoryHandler(logger *logrus.Logger, conn *utils.PostgresConn, config *cfg.Config, uc usecase.HistoryUsecase) *proxyHistory {
	return &proxyHistory{
		logger: logger,
		conn:   conn,
		config: config,
		uc:     uc,
	}
}
func (h *proxyHistory) WarnLog(format string, args ...interface{}) {
	h.logger.Warnf(format, args...)
}
func (h *proxyHistory) ErrLog(format string, args ...interface{}) {
	h.logger.Errorf(format, args...)
}
func (h *proxyHistory) InfoLog(format string, args ...interface{}) {
	h.logger.Infof(format, args...)
}
func (h *proxyHistory) HandleAllRequests(w http.ResponseWriter, _ *http.Request) {
	res, err := h.uc.GetRequests()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	respond, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(respond)
	w.WriteHeader(http.StatusOK)
}
func (h *proxyHistory) HandleRequestByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		err := errors.New("id in args must be only")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("error convert id from query arg, ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	req, err := h.uc.GetRequestByID(id)
	if err != nil {
		h.logger.Error("error from database - GetRequestByID ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalError.Error()))
		return
	}
	respond, err := json.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(respond)
	w.WriteHeader(http.StatusOK)
}
