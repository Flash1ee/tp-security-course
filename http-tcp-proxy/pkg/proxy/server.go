package proxy

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"http-tcp-proxy/cfg"
)

type Server struct {
	listener net.Listener
}

func NewServer(network string, addr string) (*Server, error) {
	if network != tcp {
		return nil, errors.New("not supported network, only tcp")
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not listen address: %v; error: %v\n", addr, err))
	}

	s := &Server{
		listener: listener,
	}

	return s, nil
}
func (s *Server) Read(conn net.Conn) ([]byte, error) {
	var err error
	receiveData := make([]byte, 0, bufferSize)
	data := make([]byte, 0, bufferSize)

	for err != io.EOF {
		n, err := conn.Read(data)
		if err != nil {
			break
		}
		data = data[:n]
		receiveData = append(receiveData, data...)
	}
	if err != nil && err != io.EOF {
		log.Printf("error on read data from connection; err = %v\n", err)
		return nil, err
	}

	return receiveData, nil
}
func (s *Server) HandleConn(conn net.Conn) {
	var err error
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("can not close listening connection, error: %v\n", err)
		}
	}(conn)

	rcvData, err := s.Read(conn)
	if err != nil {
		return
	}

	reqAddr, respData := getRequestAddress(rcvData)
}

func getRequestAddress(data []byte) (string, []byte) {

}
func (s *Server) Serve() {
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("can not close listening connection, error: %v\n", err)
		}
	}(s.listener)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatalf("can not accept connection, error: %v\n", err)
		}

		go s.HandleConn(conn) // close conn in s.Serve

	}
}
func Run(config *cfg.Config) {
	s, err := NewServer("tcp", config.BindAddr)
	if err != nil {
		log.Fatal(err)
	}
	s.Serve()
}
