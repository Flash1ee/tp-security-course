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
	p.logger.Infof("request info: %v\n", r)
	if r.URL.Scheme != SCHEME_HTTP {
		w.WriteHeader(http.StatusBadRequest)
		p.logger.Infof("unsupported protocol scheme = %v\n", r.URL.Scheme)
		return
	}

	client := httpClient()

	delHeaders(r.Header)
	r.RequestURI = ""

	resp, err := client.Do(r)
	if err != nil {
		p.logger.Errorf("ServeHTTP: %v\n", err)

		if resp.StatusCode == http.StatusNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			p.logger.Errorf("close error: %v\n", err)
		}
	}(resp.Body)

	delHeaders(resp.Header)
	copyHeader(w.Header(), resp.Header)

	w.WriteHeader(resp.StatusCode)

	w.Header()
	if _, err = io.Copy(w, resp.Body); err != nil {
		p.logger.Errorf("copy error: %v", err)
		return
	}
}
