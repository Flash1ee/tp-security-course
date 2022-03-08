package proxy

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/delivery/proxy_handler"
	"http-proxy/pkg/utils"
)

type Proxy struct {
	logger  *logrus.Logger
	conn    *utils.PostgresConn
	address string
}

func (p *Proxy) Start(config *cfg.Config) error {
	h := proxy_handler.NewProxyHandler(p.logger, p.conn, config)

	return http.ListenAndServe(p.address, h)
}

func New(log *logrus.Logger, conn *utils.PostgresConn, addr string) *Proxy {
	return &Proxy{
		logger:  log,
		conn:    conn,
		address: addr,
	}

}
