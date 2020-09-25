package clients

import (
	"bytes"
	"crypto/tls"
	"errors"
	. "github.com/ajs124/WoTest/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	auth                 AuthenticationData
	AllowSelfSignedCerts bool
}

func (c *HttpClient) Recv(u url.URL) ([]byte, MeasurementData, error) {
	var md MeasurementData
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.AllowSelfSignedCerts},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, md, err
	}
	if c.auth.Scheme == AuthBasic {
		req.SetBasicAuth(c.auth.Data["user"], c.auth.Data["password"])
	}

	md.StartTime = time.Now().UTC()
	resp, err := cl.Do(req)
	if err != nil {
		return nil, md, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	md.StopTime = time.Now().UTC()
	md.Size = uint(len(body))
	cl.CloseIdleConnections()
	return body, md, err
}

func (c *HttpClient) Send(u url.URL, contentType string, data []byte) (MeasurementData, error) {
	var md MeasurementData
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.AllowSelfSignedCerts},
	}
	cl := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", u.String(), bytes.NewReader(data))
	if err != nil {
		return md, err
	}
	req.Header.Set("Content-Type", contentType)
	if c.auth.Scheme == AuthBasic {
		req.SetBasicAuth(c.auth.Data["user"], c.auth.Data["password"])
	}

	md.Size = uint(len(data))
	md.StartTime = time.Now().UTC()
	_, err = cl.Do(req)
	md.StopTime = time.Now().UTC()
	cl.CloseIdleConnections()
	return md, err
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

func (c *HttpClient) Disconnect() error {
	return nil
}
