package clients

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct{}

func (c *HttpClient) recv(u url.URL) ([]byte, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

func (c *HttpClient) send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	_, err := http.Post(u.String(), contentType, dataReader)
	return err
}

/*
	connect(uri string) (error, response []byte) // only needed by some protocols(?)
	read(uri string) (error, []byte)
	write(uri string, data []byte) error
	subscribe(uri string) error // FIXME: needs to return some kind of pipe
	unsubscribe(uri string) error
*/
