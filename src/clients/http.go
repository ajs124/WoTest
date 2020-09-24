package clients

import (
	"bytes"
	"crypto/tls"
	"errors"
	. "github.com/ajs124/WoTest/config"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClient struct {
	auth                 AuthenticationData
	allowSelfSignedCerts bool
}

func (c *HttpClient) Recv(u url.URL) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.allowSelfSignedCerts},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	if c.auth.Scheme == AuthBasic {
		req.SetBasicAuth(c.auth.Data["user"], c.auth.Data["password"])
	}

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *HttpClient) Send(u url.URL, contentType string, data []byte) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.allowSelfSignedCerts},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	if c.auth.Scheme == AuthBasic {
		req.SetBasicAuth(c.auth.Data["user"], c.auth.Data["password"])
	}

	_, err = cl.Do(req)
	return err
}

func (c *HttpClient) Connect(u url.URL) ([]byte, error) {
	return []byte{}, nil
}

func (c *HttpClient) Subscribe(u url.URL) error {
	return errors.New("not implemented")
}

func (c *HttpClient) Unsubscribe(u url.URL) error {
	return errors.New("not implemented")
}
