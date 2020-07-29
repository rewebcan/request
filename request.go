package request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type Option struct {
	Headers map[string]string
	Cookies []*http.Cookie
}

type Client interface {
	GET(url string, opt ...Option) ([]byte, error)
	POST(url string, payload []byte, opt ...Option) ([]byte, error)
	DO(method, url string, payload []byte, opt ...Option) ([]byte, error)
}

type client struct {
	client *http.Client
}

func New() Client {
	return &client{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *client) GET(url string, opt ...Option) ([]byte, error) {
	return doRequest(c.client, "GET", url, nil, opt...)
}

func (c *client) POST(url string, payload []byte, opt ...Option) ([]byte, error) {
	return doRequest(c.client, "POST", url, payload, opt...)
}

func (c *client) DO(method, url string, payload []byte, opt ...Option) ([]byte, error) {
	return doRequest(c.client, method, url, payload, opt...)
}

func createRequest(method string, url string, payload []byte, option ...Option) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	for _, o := range option {
		for k, v := range o.Headers {
			req.Header.Set(k, v)
		}
		for _, cookie := range o.Cookies {
			req.AddCookie(cookie)
		}
	}
	return req, nil
}

func doRequest(client *http.Client, method string, url string, payload []byte, option ...Option) ([]byte, error) {
	req, err := createRequest(method, url, payload, option...)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
