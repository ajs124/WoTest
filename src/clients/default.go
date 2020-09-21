package clients

import "net/url"

type Client interface {
	connect(url url.URL) (error, response []byte) // only needed by some protocols(?)
	recv(url url.URL) (error, []byte)
	send(url url.URL, data []byte) error
	subscribe(url url.URL) error // FIXME: needs to return some kind of pipe
	unsubscribe(url url.URL) error
}
