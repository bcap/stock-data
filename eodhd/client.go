package eodhd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bcap/stock-data/jq"
	"github.com/itchyny/gojq"
)

type Client struct {
	apiKey  string
	baseURL string
	client  http.Client
}

type ClientOption func(*Client)

func WithBaseURL(baseURL string) ClientOption {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func NewClient(apiKey string, options ...ClientOption) *Client {
	c := Client{
		baseURL: "https://eodhd.com/",
		apiKey:  apiKey,
		client:  http.Client{Transport: &defaultTransport},
	}
	for _, opt := range options {
		opt(&c)
	}
	return &c
}

var noopNormalizer = jq.Script(".")

func (c *Client) process(ctx context.Context, apiPath string, normalizer *gojq.Code, parsedHook jq.ParsedHook) ([]byte, error) {
	data, err := c.get(ctx, apiPath)
	if err != nil {
		return nil, err
	}
	if normalizer == nil {
		normalizer = noopNormalizer
	}
	data, err = jq.Run(normalizer, parsedHook, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) get(ctx context.Context, apiPath string) ([]byte, error) {
	url := c.baseURL + apiPath
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("got status code %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
