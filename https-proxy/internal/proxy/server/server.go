package proxy

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/delivery/proxy_handler"
)

type Proxy struct {
	logger  *logrus.Logger
	address string
}

func (p *Proxy) Start(config *cfg.Config) error {
	h := proxy_handler.NewProxyHandler(p.logger, config)

	return http.ListenAndServe(p.address, h)
}

func New(log *logrus.Logger, addr string) *Proxy {
	return &Proxy{
		logger:  log,
		address: addr,
	}

}
