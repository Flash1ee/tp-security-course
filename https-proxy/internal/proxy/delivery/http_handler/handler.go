package http_handler

import "github.com/sirupsen/logrus"

type httpHandler struct {
	logger *logrus.Logger
}

func NewHttpHandler(logger *logrus.Logger) *httpHandler {
	return &httpHandler{
		logger: logger,
	}
}
func (h *httpHandler) Handle() {

}
