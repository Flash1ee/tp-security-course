package proxy

import "fmt"

var (
	responseOk  = []byte("HTTP/1.1 200 OK\r\n")
	responseBad = []byte("HTTP/1.1 400 BadRequest\r\n")

	contentType = []byte("Content-Type: text/plain\r\n")
	serverName  = "http-tcp-proxy"
	serverInfo  = []byte(fmt.Sprintf("Server: %s\r\nOwner: Flash1ee\r\n\r\n", serverName))
)

const (
	tcp        = "tcp"
	bufferSize = 1024
)
