package eodhd

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

type transport struct {
	http.RoundTripper
	limitter *rate.Limiter
}

func newTransport(limitter *rate.Limiter) *transport {
	return &transport{
		RoundTripper: http.DefaultTransport,
		limitter:     limitter,
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.limitter != nil {
		t.limitter.Wait(req.Context())
	}

	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		log.Printf("%s %s -> %s", req.Method, req.URL, err)
		return resp, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%s %s -> %s ! %s", req.Method, req.URL, resp.Status, err)
		return resp, err
	}

	log.Printf("%s %s -> %s %d bytes", req.Method, req.URL, resp.Status, len(body))

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	return resp, err
}
