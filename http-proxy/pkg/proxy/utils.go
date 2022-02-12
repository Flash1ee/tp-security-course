package proxy

import (
	"net/http"
)

func httpClient() *http.Client {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		//Transport: &http.Transport{
		//	MaxIdleConnsPerHost: 1,
		//},
		Timeout: TIMEOUT,
	}

	return client
}

// https://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
}

func delHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
