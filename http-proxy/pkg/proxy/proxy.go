package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type Proxy struct {
	logger  *logrus.Logger
	address string
}

func New(log *logrus.Logger, addr string) *Proxy {
	return &Proxy{
		logger:  log,
		address: addr,
	}

}
func (p *Proxy) ListenAndServe() error {
	l, err := net.Listen("tcp", p.address)
	if err != nil {
		p.logger.Errorf("net.Listen fail: %v", err)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			p.logger.Warnf("l.Accept fail: %v", err)
			continue
		}

		go p.Handle(conn)
	}
}
func (p *Proxy) Handle(conn net.Conn) {
	defer conn.Close()

	r, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		p.logger.Warnf("http.ReadRequest fail", "err", err)
		return
	}
	if r.URL.Scheme != SCHEME_HTTP && r.Method != http.MethodConnect {
		p.logger.Infof("unsupported protocol scheme = %v\n", r.URL.Scheme)
		conn.Write([]byte(fmt.Sprintf("unsupported protocol scheme = %s\n", r.URL.Scheme)))
	}
	_, err = net.LookupIP(strings.Split(r.Host, ":")[0])
	if err != nil {
		conn.Write(concatByteRespond(responseBad, []byte(fmt.Sprintf("can not resolve host : %v\r\n", r.Host))))
		p.logger.Warnf("can not resolve host: %s", r.Host)
		return
	}
	if r.Method == http.MethodConnect {
		if r.URL.Port() == "" {
			r.URL.Host = fmt.Sprintf("%s:%d", r.URL.Host, 443)
		}
	} else {
		if r.URL.Port() == "" {
			r.URL.Host = fmt.Sprintf("%s:%d", r.URL.Host, 80)
		}
	}

	r.RequestURI = ""
	remoteConn, err := net.Dial(TCP, r.URL.Host)
	if err != nil {
		p.logger.Warnf("dial remote fail", "err", err, "addr", r.URL.Host)
		return
	}

	if r.Method == http.MethodConnect {
		_, err := conn.Write(responseConnectionEstablished)
		if err != nil {
			p.logger.Warnf("https resopnse 200 fail", "err", err)
			return
		}
	} else {
		err := r.Write(remoteConn)
		if err != nil {
			p.logger.Warnf("remote write fail", "err", err)
			return
		}
	}

	go io.Copy(conn, remoteConn)
	io.Copy(remoteConn, conn)

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			p.logger.Errorf("close error: %v\n", err)
		}
	}(r.Body)
}
