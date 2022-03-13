package http_usecase

import (
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/client"
	"http-proxy/internal/proxy/repository"
	"http-proxy/pkg/utils"
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
	//
	logger *logrus.Logger
	repo   *repository.ProxyRepository
}

func NewHttpUsecase(resp http.ResponseWriter, req *http.Request, _ *cfg.Config, logger *logrus.Logger, conn *utils.PostgresConn) (*HttpUsecase, error) {
	uc := &HttpUsecase{
		clientResponse: resp,
		clientRequest:  req,
		logger:         logger,
		repo:           repository.NewProxyRepository(conn),
	}
	if conn == nil {
		return nil, InvalidArg
	}
	return uc, nil
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

	dump, err := httputil.DumpRequest(u.clientRequest, true)
	if err != nil {
		return DumpError
	}
	req := repository.FormRequestData(u.clientRequest, dump)
	reqID, err := u.repo.InsertRequest(req)
	if err != nil {
		return err
	}
	if err = u.doRequest(); err != nil {
		return err
	}

	body, err := u.sendResponse()
	if err != nil {
		return err
	}
	resp := repository.FormResponseData(u.proxyResponse, body)
	if resp == nil {
		return LogicError
	}

	return u.repo.InsertResponse(reqID, resp)
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
func (u *HttpUsecase) sendResponse() (string, error) {
	if u.clientResponse == nil {
		//u.logger.Warnf("nil response - error")
		return "", NilResponse
	}
	defer u.proxyResponse.Body.Close()

	client.PrepareResponse(u.proxyResponse, u.clientResponse)
	u.clientResponse.WriteHeader(u.proxyResponse.StatusCode)

	if _, err := io.Copy(u.clientResponse, u.proxyResponse.Body); err != nil {
		return "", err
	}

	rawResp, err := httputil.DumpResponse(u.proxyResponse, true)
	if err != nil {
		return "", err
	}
	return string(rawResp), nil

}
