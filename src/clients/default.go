package clients

import (
	"net/url"
	"time"
)

type MeasurementData struct {
	// possible timings for connecting, tls setup, sending, waiting and receiving?
	// coap and mqtt don't really "connect" per request, do it?
	// and the libraries probably don't give us that data, sooooo
	StartTime time.Time
	StopTime  time.Time
	Size      uint // in bytes
}

type Client interface {
	Connect(u url.URL) ([]byte, error)
	Recv(u url.URL) ([]byte, MeasurementData, error)
	Send(u url.URL, contentType string, data []byte) (MeasurementData, error)
	Subscribe(u url.URL) error // FIXME: needs to return some kind of pipe
	Unsubscribe(u url.URL) error
	Disconnect() error
}
