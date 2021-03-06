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

// THttpClientOptions an optional configurations for THttpClient.
type THttpClientOptions struct {
	Client *http.Client
	Header http.Header
}

// THttpClientFactory a factory of THttpClient.
type THttpClientFactory struct {
	url     string
	options THttpClientOptions
}

// THttpClient a http client implementation for TTransport.
type THttpClient struct {
	client *http.Client
	url    *url.URL
	header http.Header

	request  *bytes.Buffer
	response *http.Response

	cache [1]byte
}

// NewDefaultTHttpClientOptions returns new default THttpClientOptions.
func NewDefaultTHttpClientOptions() THttpClientOptions {
	return THttpClientOptions{
		Header: map[string][]string{"Content-Type": {"application/x-thrift"}},
	}
}

// NewTHttpClientFactory returns new THttpClientFactory.
func NewTHttpClientFactory(urlstr string) *THttpClientFactory {
	return NewTHttpClientFactoryWithOptions(urlstr, NewDefaultTHttpClientOptions())
}

// NewTHttpClientFactoryWithOptions returns new THttpClientFactory with options.
func NewTHttpClientFactoryWithOptions(urlstr string, options THttpClientOptions) *THttpClientFactory {
	return &THttpClientFactory{urlstr, options}
}

// GetTransport returns new THttpClient.
func (f *THttpClientFactory) GetTransport(t TTransport) (TTransport, error) {
	if t, ok := t.(*THttpClient); ok && t.url != nil {
		return NewTHttpClientWithOptions(t.url.String(), f.options)
	}
	return NewTHttpClientWithOptions(f.url, f.options)
}

// NewTHttpClient returns new THttpClient.
func NewTHttpClient(urlstr string) (*THttpClient, error) {
	return NewTHttpClientWithOptions(urlstr, NewDefaultTHttpClientOptions())
}

// NewTHttpClientWithOptions returns new THttpClient with options.
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
		header:  options.Header,
		request: bytes.NewBuffer(make([]byte, 0, 1024)),
	}, nil
}

// SetHeader set a http header.
func (c *THttpClient) SetHeader(k, v string) {
	c.header.Add(k, v)
}

// DelHeader delete a http header.
func (c *THttpClient) DelHeader(k string) {
	c.header.Del(k)
}

// Write writes v to request buffer.
func (c *THttpClient) Write(v []byte) (int, error) {
	return c.request.Write(v)
}

// WriteByte writes v to request buffer.
func (c *THttpClient) WriteByte(v byte) error {
	return c.request.WriteByte(v)
}

// Read reads v from response body.
func (c *THttpClient) Read(v []byte) (int, error) {
	n, err := c.response.Body.Read(v)
	return n, NewTTransportExceptionFromError(err)
}

// ReadByte reads next one byte from response body.
func (c *THttpClient) ReadByte() (byte, error) {
	_, err := c.Read(c.cache[:])
	return c.cache[0], err
}

// Flush sends request to server.
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
	if c.response, err = c.client.Do(req); err != nil {
		return NewTTransportExceptionFromError(err)
	}
	if c.response.StatusCode != http.StatusOK {
		status := c.response.StatusCode
		c.closeResponse()
		return NewTTransportException(TTransportErrorUnknown, fmt.Sprintf("HTTP status code: %v", status))
	}
	return
}

func (c *THttpClient) closeResponse() (err error) {
	if c.response != nil {
		io.Copy(ioutil.Discard, c.response.Body)
		err = c.response.Body.Close()
		c.response = nil
	}
	return err
}
