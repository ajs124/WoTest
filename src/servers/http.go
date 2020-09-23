package servers

import (
	"errors"
	"net"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
	mux    *http.ServeMux
	tls    bool
}

func (s *HttpServer) Listen(address string, tls bool, certFile, keyFile string) error {
	serveMux := http.NewServeMux()
	s.mux = serveMux
	server := &http.Server{
		Addr:           address,
		Handler:        serveMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.server = server

	if address == "" {
		address = ":http"
	}
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	if tls {
		go server.ServeTLS(ln, certFile, keyFile)
	} else {
		go server.Serve(ln) // FIXME: handle errors? but we need to `go`
	}
	return nil
}

func (s *HttpServer) SetHandleFunc(f func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc("/", f)
}

func (s *HttpServer) Send(path string, contents []byte) error {
	return errors.New("not implemented")
}

func (s *HttpServer) Stop() error {
	return s.server.Close()
}
