package clients

import (
	"bytes"
	"context"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/rs/zerolog/log"

	"io/ioutil"
	"net"
	"net/url"
	"time"
)

type CoapClient struct {
	co net.Conn
}

func (c *CoapClient) connect(u url.URL) {
	co, err := udp.Dial(u.Host + ":" + u.Port())
	c.co = co
	// for dtls
	// co, err := dtls.Dial("localhost:5688", &dtls.Config{...}))
	if err != nil {
		log.Err(err).Msg("Coap: Error dialing.")
	}
}

func (c *CoapClient) recv(u url.URL) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := c.co.Get(ctx, u.Path) // maybe EscapedPath? who knows…
	if err != nil {
		log.Err(err).Msg("Coap: cannot get response")
	}
	return resp, err
}

func (c *CoapClient) send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	c.co.Post(ctx, u.Path)
	// _, err := http.Post(uri, contentType, dataReader)
	return err
}
