package proxy

import "time"

var (
	responseConnectionEstablished = []byte("HTTP/1.1 200 Connection established\r\n\r\n")
	responseBad                   = []byte("HTTP/1.1 400 BadRequest\r\n")
)

const (
	SCHEME_HTTP  = "http"
	SCHEME_HTTPS = "https"
	TCP          = "tcp"
	TIMEOUT      = 5 * time.Second
	MaxHostConn  = 20
	bufferSize   = 1024
)
