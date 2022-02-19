package proxy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

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
func (s *Server) Read(conn net.Conn) (*http.Request, error) {
	//receiveData := make([]byte, 0, bufferSize)
	//data := make([]byte, bufferSize)
	req, err := http.ReadRequest(bufio.NewReader(conn))
	return req, err
	//if req.URL.Hostname() == "" {
	//	respondOk(conn)
	//}
	//if err != nil {
	//	log.Printf("can not read request; err = %v\n", err)
	//	return nil, err
	//}
	//for err != io.EOF {
	//	n, err := conn.Read(data)
	//	if err != nil {
	//		if err == io.EOF {
	//			break
	//		}
	//	}
	//	data = data[:n]
	//	receiveData = append(receiveData, data...)
	//}
	//if err != nil && err != io.EOF {
	//	log.Printf("error on read data from connection; err = %v\n", err)
	//	return nil, err
	//}
	//
	//return receiveData, nil
}

func respond(conn net.Conn, arg []byte) {
	_, err := conn.Write(arg)
	if err != nil {
		log.Printf("write respond error; err: %v\n", err)
	}
	return
}
func (s *Server) HandleConn(conn net.Conn) {
	var err error
	defer func(conn net.Conn) {
		err = conn.Close()
		if err != nil {
			log.Fatalf("can not close listening connection, error: %v\n", err)
		}
	}(conn)

	req, err := s.Read(conn)
	if err != nil {
		log.Printf(fmt.Sprintf("http.ReadRequest fail; err: %v\n", err))
		conn.Write([]byte(fmt.Sprintf("Http request error: %d", http.StatusBadRequest)))
		return
	}
	if req.URL.Scheme != "http" || req.Method == http.MethodConnect {
		log.Printf(fmt.Sprintf("unsupported protocol scheme = %s\n", req.URL.Scheme))
		respond(conn, concatByteRespond(responseBad, contentType, []byte(fmt.Sprintf("unsupported protocol scheme = %s\n", req.URL.Scheme))))
		return
	}
	_, err = net.LookupIP(req.Host)
	if err != nil {
		respond(conn, concatByteRespond(responseBad, contentType, []byte(fmt.Sprintf("Could not resolve host: %s\r\n", req.Host))))
		return
	}

	if req.Method == http.MethodConnect {
		if req.URL.Port() == "" {
			req.URL.Host = fmt.Sprintf("%s:%d", req.URL.Host, 443)
		}
	} else {
		if req.URL.Port() == "" {
			req.URL.Host = fmt.Sprintf("%s:%d", req.URL.Host, 80)
		}
	}
	if req.URL.Hostname() == "" {
		respond(conn, concatByteRespond(responseOk, contentType, serverInfo))
		return
	}
	req.RequestURI = ""
	log.Printf("newRequest, url: %s\n", req.URL.String())
	remoteConn, err := net.Dial(tcp, req.URL.Host)
	if err != nil {
		log.Printf("dial remote fail err: %s addr: %s\n", err, req.URL.Host)
		return
	}
	defer remoteConn.Close()
	if req.Method == http.MethodConnect {
		// response ok
		_, err = conn.Write(responseConnectionEstablished)
		if err != nil {
			log.Printf("https resopnse 200 fail err: %s\n", err)
			return
		}
		return
	}
	err = req.Write(remoteConn)
	if err != nil {
		log.Printf("remote write line fail err: %s\n", err)
		return
	}
	io.Copy(conn, remoteConn)
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
