package client

import "net/http"

// https://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Handle-Authenticate",
	"Handle-Authorization",
	"Handle-Connection",
}

func PrepareResponse(proxyResp *http.Response, clientResp http.ResponseWriter) {
	copyHeader(clientResp.Header(), proxyResp.Header)
}
func PrepareRequest(req *http.Request, proxyRequest *http.Request) {
	copyHeader(proxyRequest.Header, req.Header)
	delHeaders(proxyRequest.Header)
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
