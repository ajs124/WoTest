package clients

// always pass a URI
type Client interface {
	connect(uri string) (error, response []byte) // only needed by some protocols(?)
	read(uri string) (error, []byte)
	write(uri string, data []byte) error
	observe(uri string) error // FIXME: needs to return some kind of pipe
	unobserve(uri string) error
	subscribe(uri string) error // FIXME: needs to return some kind of pipe
	unsubscribe(uri string) error
}
