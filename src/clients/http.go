package clients

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct{}

func (c *HttpClient) Recv(u url.URL) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return bytes, err
}

func (c *HttpClient) Send(u url.URL, contentType string, data []byte) error {
	dataReader := bytes.NewReader(data)
	_, err := http.Post(u.String(), contentType, dataReader)
	return err
}

func (c *HttpClient) Connect(u url.URL) ([]byte, error) {
	return []byte{}, nil
}

func (c *HttpClient) Subscribe(u url.URL) error {
	return nil
}

func (c *HttpClient) Unsubscribe(u url.URL) error {
	return nil
}
