package servers

type Server interface {
	Listen(address string) error
	// SetHandleFunc(f func(w, r interface{})) // why am I even trying? -.-
	Send(path string, contents []byte) error
	Stop() error
}
