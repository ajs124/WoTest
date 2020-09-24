package servers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

type CoapServer struct {
	router *mux.Router
	server *udp.Server
	conn   *net.UDPConn
}

func logRequest(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Debug().Str("client address", w.Client().RemoteAddr().String()).Str("request", r.String()).Msg("coap server received")
		next.ServeCOAP(w, r)
	})
}

func getPath(opts message.Options) string {
	path, err := opts.Path()
	if err != nil {
		log.Printf("cannot get path: %v", err)
		return ""
	}
	return path
}

func sendResponse(cc mux.Client, token []byte, obs int64, payload []byte, contentFormat message.MediaType) error {
	m := message.Message{
		Code:    codes.Content,
		Token:   token,
		Context: cc.Context(),
		Body:    bytes.NewReader(payload),
	}
	var opts message.Options
	var buf []byte
	opts, n, err := opts.SetContentFormat(buf, contentFormat)
	if err == message.ErrTooSmall {
		buf = append(buf, make([]byte, n)...)
		opts, n, err = opts.SetContentFormat(buf, contentFormat)
	}
	if err != nil {
		return fmt.Errorf("cannot set content format to response: %w", err)
	}
	if obs >= 0 {
		opts, n, err = opts.SetObserve(buf, uint32(obs))
		if err == message.ErrTooSmall {
			buf = append(buf, make([]byte, n)...)
			opts, n, err = opts.SetObserve(buf, uint32(obs))
		}
		if err != nil {
			return fmt.Errorf("cannot set options to response: %w", err)
		}
	}
	m.Options = opts
	return cc.WriteMessage(&m)
}

func (s *CoapServer) ServeFile(w mux.ResponseWriter, r *mux.Message, filePath string) {
	log.Printf("Got message path=%v: %+v from %v", getPath(r.Options), r, w.Client().RemoteAddr())
	// obs, err := r.Options.Observe()
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error().Err(err).Msg("error reading file")
	}
	/* switch {
	case r.Code == codes.GET && err == nil && obs == 0:
		// go periodicTransmitter(w.Client(), r.Token)
	case r.Code == codes.GET:
		subded := time.Now() */
	err = sendResponse(w.Client(), r.Token, -1, fileContents, message.AppJSON)
	if err != nil {
		log.Printf("Error on transmitter: %v", err)
	}
	// }
}

func (s *CoapServer) Listen(address string) error {
	r := mux.NewRouter()
	r.Use(logRequest)
	s.router = r

	// FIXME: support networks other than UDP?
	l, err := net.NewListenUDP("udp", address)
	if err != nil {
		log.Error().Err(err).Msg("error listening on udp")
		return err
	}
	s.conn = l
	server := udp.NewServer(udp.WithMux(s.router))
	s.server = server
	go server.Serve(l) // FIXME: handle errors? but we need to `go`
	return err
}

func (s *CoapServer) SetHandleFunc(f func(mux.ResponseWriter, *mux.Message)) {
	s.router.DefaultHandleFunc(f)
}

func (s *CoapServer) Send(path string, contents []byte) error {
	return errors.New("not implemented")
}

func (s *CoapServer) Stop() error {
	s.server.Stop()
	s.conn.Close()
	return nil
}
