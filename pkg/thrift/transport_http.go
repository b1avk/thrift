package thrift

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type THttpClientOptions struct {
	Client *http.Client
}

type THttpClientFactory struct {
	url     string
	options THttpClientOptions
}

type THttpClient struct {
	client *http.Client
	url    *url.URL
	header http.Header

	request  *bytes.Buffer
	response *http.Response

	cache [1]byte
}

func NewTHttpClientFactory(urlstr string) *THttpClientFactory {
	return NewTHttpClientFactoryWithOptions(urlstr, THttpClientOptions{})
}

func NewTHttpClientFactoryWithOptions(urlstr string, options THttpClientOptions) *THttpClientFactory {
	return &THttpClientFactory{urlstr, options}
}

func (f *THttpClientFactory) GetTransport(t TTransport) (TTransport, error) {
	if t, ok := t.(*THttpClient); ok && t.url != nil {
		return NewTHttpClientWithOptions(t.url.String(), f.options)
	}
	return NewTHttpClientWithOptions(f.url, f.options)
}

func NewTHttpClientWithOptions(urlstr string, options THttpClientOptions) (*THttpClient, error) {
	parsedURL, err := url.Parse(urlstr)
	if err != nil {
		return nil, err
	}
	client := options.Client
	if client == nil {
		client = http.DefaultClient
	}
	return &THttpClient{
		client:  client,
		url:     parsedURL,
		header:  map[string][]string{"Content-Type": {"application/x-thrift"}},
		request: bytes.NewBuffer(make([]byte, 0, 1024)),
	}, nil
}

func NewTHttpClient(urlstr string) (*THttpClient, error) {
	return NewTHttpClientWithOptions(urlstr, THttpClientOptions{})
}

func (c *THttpClient) SetHeader(k, v string) {
	c.header.Add(k, v)
}

func (c *THttpClient) DelHeader(k string) {
	c.header.Del(k)
}

func (c *THttpClient) Write(v []byte) (int, error) {
	return c.request.Write(v)
}

func (c *THttpClient) WriteByte(v byte) error {
	return c.request.WriteByte(v)
}

func (c *THttpClient) Read(v []byte) (int, error) {
	n, err := c.response.Body.Read(v)
	return n, NewTTransportExceptionFromError(err)
}

func (c *THttpClient) ReadByte() (byte, error) {
	_, err := c.Read(c.cache[:])
	return c.cache[0], err
}

func (c *THttpClient) Flush(ctx context.Context) (err error) {
	c.closeResponse()
	buf := c.request
	c.request = new(bytes.Buffer)
	req, err := http.NewRequest("POST", c.url.String(), buf)
	if err != nil {
		return NewTTransportExceptionFromError(err)
	}
	req.Header = c.header
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	c.response, err = c.client.Do(req)
	if c.response.StatusCode != http.StatusOK {
		c.closeResponse()
		return NewTTransportException(TTransportErrorUnknown, fmt.Sprintf("HTTP status code: %v", c.response.StatusCode))
	}
	return
}

func (c *THttpClient) closeResponse() (err error) {
	if c.response != nil {
		io.Copy(ioutil.Discard, c.response.Body)
		err = c.response.Body.Close()
	}
	c.response = nil
	return err
}
