package proxy

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Proxy struct {
	logger *logrus.Logger
}

func New(log *logrus.Logger) *Proxy {
	return &Proxy{
		logger: log,
	}

}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := httpClient()

	delHeaders(r.Header)
	r.RequestURI = ""

	resp, err := client.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.logger.Errorf("ServeHTTP: %v\n", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			p.logger.Errorf("close error: %v\n", err)
		}
	}(resp.Body)

	delHeaders(resp.Header)
	w.WriteHeader(resp.StatusCode)

	copyHeader(w.Header(), resp.Header)

	if _, err = io.Copy(w, resp.Body); err != nil {
		p.logger.Errorf("copy error: %v", err)
	}
}
