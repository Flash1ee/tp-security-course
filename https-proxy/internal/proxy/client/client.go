package client

import (
	"net/http"
	"time"
)

const (
	TIMEOUT = 10 * time.Second
)

func HttpClient() *http.Client {
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
