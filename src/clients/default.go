package clients

import "net/url"

type Client interface {
	Connect(u url.URL) ([]byte, error)
	Recv(u url.URL) ([]byte, error)
	Send(u url.URL, contentType string, data []byte) error
	Subscribe(u url.URL) error // FIXME: needs to return some kind of pipe
	Unsubscribe(u url.URL) error
}
