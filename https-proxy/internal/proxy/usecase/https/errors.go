package https_usecase

import "errors"

var (
	ReadClientReqError = errors.New("read client request error")
	DumpError          = errors.New("can not get dump")
	LoadCertError      = errors.New("load server certificate error")
	InvalidArg         = errors.New("invalid arg")
	ProxyConnErr       = errors.New("invalid host")
)
