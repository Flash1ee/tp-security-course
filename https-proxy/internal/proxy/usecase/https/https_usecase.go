package https_usecase

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"

	"http-proxy/cfg"
	"http-proxy/internal/proxy/repository"
	"http-proxy/internal/proxy/usecase"
	"http-proxy/pkg/utils"
)

//+------+            +-----+                   +-----------+
//|client|            |proxy|                   |destination|
//+------+            +-----+                   +-----------+
//          --CONNECT-->
//         <-- Conn Established--
//                             <--TCP handshake-->
//           <--------------Tunnel---------------->

type HttpsUsecase struct {
	clientResponse http.ResponseWriter
	clientRequest  *http.Request
	proxyResponse  *http.Response
	proxyRequest   *http.Request
	//
	tlsConfig  *tls.Config
	proxyConn  *tls.Conn
	clientConn net.Conn
	//
	logger *logrus.Logger
	repo   *repository.ProxyRepository
}

func (u *HttpsUsecase) Close() {
	u.clientConn.Close()
	u.proxyConn.Close()

	if u.clientRequest != nil {
		u.clientRequest.Body.Close()
	}
	if u.proxyRequest != nil {
		u.proxyRequest.Body.Close()
	}
}
func NewHttpsUsecase(resp http.ResponseWriter, req *http.Request, config *cfg.Config, logger *logrus.Logger, conn *utils.PostgresConn) (*HttpsUsecase, error) {
	uc := &HttpsUsecase{
		clientResponse: resp,
		clientRequest:  req,
		logger:         logger,
		repo:           repository.NewProxyRepository(conn),
	}
	if conn == nil {
		return nil, InvalidArg
	}
	var err error

	if err = uc.SetupCerts(config.CertKeyPath); err != nil {
		return nil, err
	}

	if uc.proxyConn, err = setupProxyConn(req, uc.tlsConfig); err != nil {
		return nil, ProxyConnErr
	}

	if uc.clientConn, err = setupClientConn(resp, uc.tlsConfig); err != nil {
		return nil, err
	}

	return uc, nil

}
func (uc *HttpsUsecase) SetupCerts(certKey string) error {
	pwd, _ := os.Getwd()
	rootCerts := pwd + "/cert/"
	hostCerts := pwd + "/certs/"

	hostName, err := url.Parse(uc.clientRequest.RequestURI)
	if err != nil {
		return err
	}
	curHostCert := fmt.Sprintf("%s%s%s", hostCerts, hostName.Scheme, ".crt")
	if _, err = os.Stat(curHostCert); os.IsNotExist(err) {
		err = utils.GenHostCert(rootCerts, "gen_cert.sh", hostName.Scheme, hostCerts)
		if err != nil {
			return err
		}
	}
	serverCert, err := tls.LoadX509KeyPair(curHostCert, pwd+certKey)
	if err != nil {
		uc.logger.Errorf(err.Error())
		return LoadCertError
	}
	uc.tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ServerName:   hostName.Scheme,
	}

	return nil
}
func setupClientConn(resp http.ResponseWriter, cfg *tls.Config) (net.Conn, error) {
	clientTcpConn, _, err := resp.(http.Hijacker).Hijack()
	if err != nil {
		return nil, err
	}
	if _, err = clientTcpConn.Write(usecase.ResponseConnectionEstablished); err != nil {
		return nil, err
	}

	clientConn := tls.Server(clientTcpConn, cfg)
	if err = clientConn.Handshake(); err != nil {
		return nil, err
	}

	return clientConn, nil
}
func setupProxyConn(req *http.Request, tlsCfg *tls.Config) (*tls.Conn, error) {
	destConn, err := tls.Dial("tcp", req.Host, tlsCfg)
	if err != nil {
		return nil, err
	}
	return destConn, nil
}

func (u *HttpsUsecase) Handle() error {
	if _, err := u.clientConn.Write(usecase.ResponseConnectionEstablished); err != nil {
		return err
	}
	reqConnect := repository.FormRequestData(u.clientRequest)
	_, err := u.repo.InsertRequest(reqConnect)
	if err != nil {
		return err
	}

	if err = u.getClientRequest(); err != nil {
		return ReadClientReqError
	}

	req := repository.FormRequestData(u.clientRequest)
	reqID, err := u.repo.InsertRequest(req)
	if err != nil {
		return err
	}

	dump, err := httputil.DumpRequest(u.clientRequest, true)
	if err != nil {
		return DumpError
	}
	u.logger.Infof("https client request dump: %v\n", string(dump))
	if err = u.doServerRequest(); err != nil {
		return err
	}
	body, err := u.sendResponse(u.proxyResponse)
	if err != nil {
		return err
	}

	resp := repository.FormResponseData(u.proxyResponse, body)
	if resp == nil {
		return LogicError
	}
	if err = u.repo.InsertResponse(reqID, resp); err != nil {
		return err
	}
	return nil
}
func (u *HttpsUsecase) getClientRequest() error {
	var err error
	reader := bufio.NewReader(u.clientConn)
	u.clientRequest, err = http.ReadRequest(reader)
	if err != nil {
		return err
	}

	return nil
}
func (u *HttpsUsecase) doServerRequest() error {
	if u.clientRequest == nil {
		return InvalidArg
	}
	dump, err := httputil.DumpRequest(u.clientRequest, true)

	if err != nil {
		return err
	}

	if _, err = u.proxyConn.Write(dump); err != nil {
		return err
	}
	//@TODO блокируется для google.com/youtube.com
	u.proxyResponse, err = http.ReadResponse(bufio.NewReader(u.proxyConn), u.clientRequest)

	return err
}
func (u *HttpsUsecase) sendResponse(serverResp *http.Response) (string, error) {
	rawResp, err := httputil.DumpResponse(serverResp, true)
	_, err = u.clientConn.Write(rawResp)
	if err != nil {
		return "", err
	}
	return string(rawResp), nil
}
