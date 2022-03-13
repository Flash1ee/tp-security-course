package proxy_history

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/client"
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
func (h *proxyHistory) HandleRetryRequest(w http.ResponseWriter, r *http.Request) {
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
	buffer, err := http.ReadRequest(bufio.NewReader(strings.NewReader(req.Raw)))
	if err != nil {
		h.logger.Error("error convert request from database to buffer ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalError.Error()))
		return
	}
	host, ok := req.Headers["Host"].(string)
	if !ok {
		h.logger.Error("can not convert interface to string Host ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalError.Error()))
		return
	}
	proxyReq, err := http.NewRequest(buffer.Method, req.Path, buffer.Body)
	if err != nil {
		h.logger.Error("error convert request from database to buffer ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(InternalError.Error()))
		return
	}
	proxyReq.Host = host
	proxyReq.URL.Host = host
	proxyReq.URL.Scheme = "http"
	proxyReq.URL.Opaque = ""

	utils.CopyHeader(proxyReq.Header, buffer.Header)

	checkCommandInjection := false
	if strings.Contains(req.Raw, "cat /etc/passwd;") {
		checkCommandInjection = true
	}
	var resp *http.Response
	var rawResp []byte
	if req.IsHTTPS {
		proxyReq.Host = fmt.Sprintf("%s:%s", host, "443")
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err := tls.Dial("tcp", proxyReq.Host, conf)
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		conn.Write([]byte(req.Raw))
		resp, err = http.ReadResponse(bufio.NewReader(conn), proxyReq)
		rawResp, err = httputil.DumpResponse(resp, true)

		fmt.Println(string(rawResp))
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
	} else {
		clientConn := client.HttpClient()
		resp, err = clientConn.Do(proxyReq)
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		rawResp, err = httputil.DumpResponse(resp, true)
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)

	if checkCommandInjection {
		if strings.Contains(string(rawResp), "root:") {
			h.logger.Warn(proxyReq.Host + "with command injection problem!!!")
			_, _ = fmt.Fprintf(w, "\nWARN: This request have command injection vulnerability")
			return
		}
	}
	_, _ = fmt.Fprintf(w, "\nINFO: This request haven't command injection vulnerability")

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
