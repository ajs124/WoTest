package servers

import (
	"errors"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/rs/zerolog/log"
)

type CoapServer struct {
	router *mux.Router
	server *udp.Server
}

func logRequest(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Debug().Str("client address", w.Client().RemoteAddr().String()).Str("request", r.String()).Msg("coap server received")
		next.ServeCOAP(w, r)
	})
}

func (s *CoapServer) Listen(address string) error {
	r := mux.NewRouter()
	r.Use(logRequest)
	s.router = r

	// FIXME: support networks other than UDP?
	l, err := net.NewListenUDP("udp", address)
	if err != nil {
		return err
	}
	defer l.Close()
	server := udp.NewServer(udp.WithMux(s.router))
	s.server = server
	go server.Serve(l) // FIXME: handle errors? but we need to `go`
	return err
}

func (s *CoapServer) SetHandleFunc(f func(mux.ResponseWriter, *mux.Message)) {
	s.router.HandleFunc("/", f)
}

func (s *CoapServer) Send(path string, contents []byte) error {
	return errors.New("not implemented")
}

func (s *CoapServer) Stop() error {
	// FIXME: this segfaults, why?
	// s.server.Stop()
	return nil
}
