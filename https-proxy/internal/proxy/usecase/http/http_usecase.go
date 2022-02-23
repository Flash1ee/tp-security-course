package http_usecase

import (
	"io"
	"net/http"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/client"
)

//+------+        +-----+        +-----------+
//|client|        |proxy|        |destination|
//+------+        +-----+        +-----------+
//          --Req-->
//                         --Req-->
//                         <--Res--
//          <--Res--

type HttpUsecase struct {
	clientResponse http.ResponseWriter
	clientRequest  *http.Request
	proxyResponse  *http.Response
	proxyRequest   *http.Request
	logger         *logrus.Logger
}

func NewHttpUsecase(resp http.ResponseWriter, req *http.Request, _ *cfg.Config, logger *logrus.Logger) (*HttpUsecase, error) {
	return &HttpUsecase{
		clientResponse: resp,
		clientRequest:  req,
		logger:         logger,
	}, nil
}
func (u *HttpUsecase) Close() {
	u.clientRequest.Body.Close()
	u.proxyRequest.Body.Close()
}
func (u *HttpUsecase) Handle() error {
	var err error
	u.proxyRequest, err = http.NewRequest(u.clientRequest.Method, u.clientRequest.RequestURI, u.clientRequest.Body)
	if err != nil {
		return err
	}

	if err = u.doRequest(); err != nil {
		return err
	}
	return u.sendResponse()
}
func (u *HttpUsecase) doRequest() error {
	if u.clientRequest == nil {
		//u.logger.Warnf("empty request - error")
		return NilError
	}
	var err error
	c := client.HttpClient()

	client.PrepareRequest(u.proxyRequest, u.clientRequest)
	u.proxyRequest.RequestURI = ""
	if u.proxyResponse, err = c.Do(u.proxyRequest); err != nil {
		if u.proxyResponse.StatusCode == http.StatusNotFound {
			u.clientResponse.WriteHeader(http.StatusNotFound)
		} else {
			u.clientResponse.WriteHeader(http.StatusInternalServerError)
		}
		return err
	}
	return nil

}
func (u *HttpUsecase) sendResponse() error {
	if u.clientResponse == nil {
		//u.logger.Warnf("nil response - error")
		return NilResponse
	}
	defer u.proxyResponse.Body.Close()

	client.PrepareResponse(u.proxyResponse, u.clientResponse)
	u.clientResponse.WriteHeader(u.proxyResponse.StatusCode)

	if _, err := io.Copy(u.clientResponse, u.proxyResponse.Body); err != nil {
		return err
	}

	return nil

}
