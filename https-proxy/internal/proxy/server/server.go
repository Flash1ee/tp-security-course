package proxy

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/delivery/proxy_handler"
	"http-proxy/internal/proxy/delivery/proxy_history"
	"http-proxy/internal/proxy/repository"
	"http-proxy/internal/proxy/usecase/history"
	"http-proxy/pkg/utils"
)

type Proxy struct {
	logger  *logrus.Logger
	conn    *utils.PostgresConn
	address string
}

func (p *Proxy) Start(config *cfg.Config) error {
	proxyRepo := repository.NewProxyRepository(p.conn)
	if proxyRepo == nil {
		p.logger.Fatal("Proxy repository can not init")
	}
	ucHistory := history.NewHistoryUsecase(p.logger, proxyRepo)
	if ucHistory == nil {
		p.logger.Fatal("Proxy usecase can not init")

	}

	proxyHandler := proxy_handler.NewProxyHandler(p.logger, p.conn, config)
	historyHandler := proxy_history.NewProxyHistoryHandler(p.logger, p.conn, config, ucHistory)

	reqRouter := mux.NewRouter()

	reqRouter.HandleFunc("/requests", historyHandler.HandleAllRequests)
	reqRouter.HandleFunc("/request/{id:[0-9]+}", historyHandler.HandleRequestByID)
	reqRouter.HandleFunc("/request/retry/{id:[0-9]+}", historyHandler.HandleRetryRequest)

	go http.ListenAndServe(":8081", reqRouter)
	return http.ListenAndServe(p.address, proxyHandler)
}

func New(log *logrus.Logger, conn *utils.PostgresConn, addr string) *Proxy {
	return &Proxy{
		logger:  log,
		conn:    conn,
		address: addr,
	}

}
