package clients

import (
	"context"
	"errors"
	"github.com/plgd-dev/go-coap/v2/udp"
	"github.com/plgd-dev/go-coap/v2/udp/client"
	"github.com/rs/zerolog/log"
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

func (c *CoapClient) Recv(u url.URL) ([]byte, MeasurementData, error) {
	var md MeasurementData
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	md.StartTime = time.Now().UTC()
	msg, err := c.co.Get(ctx, u.Path) // maybe EscapedPath? who knows…
	if err != nil {
		log.Err(err).Msg("Coap: cannot get response")
		return nil, md, err
	}
	resp, err := msg.ReadBody()
	md.StopTime = time.Now().UTC()
	md.Size = uint(len(resp))
	return resp, md, err
}

/* func (c *CoapClient) Send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	c.co.Post(ctx, u.Path)
	// _, err := http.Post(uri, contentType, dataReader)
	return err
} */

func (c *CoapClient) Send(u url.URL, contentType string, data []byte) (MeasurementData, error) {
	var md MeasurementData
	// dataReader := bytes.NewReader(data)
	return md, errors.New("not implemented yet")
}

func (c *CoapClient) Subscribe(u url.URL) error {
	return nil
}

func (c *CoapClient) Unsubscribe(u url.URL) error {
	return nil
}

func (c *CoapClient) Disconnect() error {
	return c.co.Close()
}
