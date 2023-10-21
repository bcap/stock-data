package eodhd

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

var defaultTransport = loggedTransport{
	RoundTripper: http.DefaultTransport,
}

type loggedTransport struct {
	http.RoundTripper
}

func (t *loggedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
