package clients

import (
	"bytes"
	"context"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/rs/zerolog/log"
	"net/http"

	"net/url"
	"time"
)

type CoapClient struct {
	co *client.ClientConn
}

func (c *CoapClient) Connect(u url.URL) ([]byte, error) {
	port := u.Port()
	if port == "" {
		port = "5683"
	}
	co, err := udp.Dial(u.Hostname() + ":" + port)
	c.co = co
	// for dtls
	// co, err := dtls.Dial("localhost:5688", &dtls.Config{...}))
	if err != nil {
		log.Err(err).Msg("Coap: Error dialing.")
	}
	return []byte{}, err
}

func (c *CoapClient) Recv(u url.URL) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msg, err := c.co.Get(ctx, u.Path) // maybe EscapedPath? who knows…
	if err != nil {
		log.Err(err).Msg("Coap: cannot get response")
		return []byte{}, err
	}
	resp, err := msg.ReadBody()
	return resp, err
}

/* func (c *CoapClient) Send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	c.co.Post(ctx, u.Path)
	// _, err := http.Post(uri, contentType, dataReader)
	return err
} */

func (c *CoapClient) Send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	_, err := http.Post(u.String(), contentType, dataReader)
	return err
}

func (c *CoapClient) Subscribe(u url.URL) error {
	return nil
}

func (c *CoapClient) Unsubscribe(u url.URL) error {
	return nil
}
