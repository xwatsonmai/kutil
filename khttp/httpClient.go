package khttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type httpClient struct {
	client *http.Client
}

type httpClientOption func(*httpClient)

func WithMaxIdleConns(maxIdleConns int) httpClientOption {
	return func(c *httpClient) {
		c.client.Transport.(*http.Transport).MaxIdleConns = maxIdleConns
	}
}

func WithMaxIdleConnsPerHost(maxIdleConnsPerHost int) httpClientOption {
	return func(c *httpClient) {
		c.client.Transport.(*http.Transport).MaxIdleConnsPerHost = maxIdleConnsPerHost
	}
}

func WithIdleConnTimeout(idleConnTimeout time.Duration) httpClientOption {
	return func(c *httpClient) {
		c.client.Transport.(*http.Transport).IdleConnTimeout = idleConnTimeout
	}
}

type httpClientOptionSet struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
}

func NewHttpClientOptionSet(options ...httpClientOption) *httpClientOptionSet {
	c := &httpClientOptionSet{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	}
	for _, option := range options {
		option(&httpClient{client: &http.Client{Transport: &http.Transport{}}})
	}
	return c
}

func NewHttpClient(optionSet *httpClientOptionSet) *httpClient {
	c := &httpClient{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        optionSet.MaxIdleConns,
				MaxIdleConnsPerHost: optionSet.MaxIdleConnsPerHost,
				IdleConnTimeout:     optionSet.IdleConnTimeout,
			},
		},
	}
	return c
}

type Result struct {
	bytes []byte
	error error
}

func (r *Result) ToStruct(v interface{}) error {
	err := json.Unmarshal(r.bytes, v)
	if err != nil {
		return err
	}
	return nil
}

func (r *Result) ToBytes() ([]byte, error) {
	if r.error != nil {
		return nil, r.error
	}
	return r.bytes, nil
}

func (r *Result) ToMap() (map[string]interface{}, error) {
	if r.error != nil {
		return nil, r.error
	}
	var m map[string]interface{}
	err := json.Unmarshal(r.bytes, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *httpClient) Get(url string) Result {
	resp, err := c.client.Get(url)
	if err != nil {
		return Result{nil, err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{nil, err}
	}

	return Result{body, nil}
}

func (c *httpClient) Post(url string, data interface{}) Result {
	var body []byte
	var err error
	switch d := data.(type) {
	case []byte:
		body = d
	case map[string]interface{}:
		body, err = json.Marshal(d)
		if err != nil {
			return Result{nil, err}
		}
	case *struct{}:
		body, err = json.Marshal(d)
		if err != nil {
			return Result{nil, err}
		}
	default:
		return Result{nil, errors.New("data类型错误")}
	}
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return Result{nil, err}
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return Result{nil, err}
	}

	return Result{body, nil}
}
