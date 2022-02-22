package http_usecase

import (
	"io"
	"net/http"

	"http-proxy/internal/proxy/client"
)

type HttpUsecase struct {
	response      http.ResponseWriter
	proxyResponse *http.Response
}

func (u *HttpUsecase) Handle(resp http.ResponseWriter, r *http.Request, req *http.Request) error {
	u.response = resp
	if err := u.doRequest(r, req); err != nil {
		return err
	}
	return u.sendResponse()
}
func (u *HttpUsecase) doRequest(req *http.Request, proxyrReq *http.Request) error {
	if req == nil {
		//u.logger.Warnf("empty request - error")
		return NilError
	}
	var err error
	c := client.HttpClient()

	client.PrepareRequest(proxyrReq, req)
	proxyrReq.RequestURI = ""
	if u.proxyResponse, err = c.Do(proxyrReq); err != nil {
		if u.proxyResponse.StatusCode == http.StatusNotFound {
			u.response.WriteHeader(http.StatusNotFound)
		} else {
			u.response.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}
	return nil

}
func (u *HttpUsecase) sendResponse() error {
	if u.response == nil {
		//u.logger.Warnf("nil response - error")
		return NilResponse

	}
	client.PrepareResponse(u.proxyResponse, u.response)
	u.response.WriteHeader(u.proxyResponse.StatusCode)

	if _, err := io.Copy(u.response, u.proxyResponse.Body); err != nil {
		return err
	}

	return nil

}
