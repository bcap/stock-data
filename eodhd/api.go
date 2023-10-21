package eodhd

import (
	"context"
	"fmt"
	"time"

	"github.com/bcap/stock-data/config"
	"github.com/bcap/stock-data/jq"
)

var normalizeExchanges = jq.LoadScript("normalizers/exchanges.jq")
var normalizeFundamentals = jq.LoadScript("normalizers/fundamentals.jq")
var normalizeHistoricalIntraday = jq.LoadScript("normalizers/historical-intraday.jq")

func (c *Client) Exchanges(ctx context.Context) ([]byte, error) {
	apiPath := fmt.Sprintf("api/exchanges-list/?api_token=%s", c.apiKey)
	return c.process(ctx, apiPath, normalizeExchanges, nil)
}

func (c *Client) Fundamentals(ctx context.Context, ticker config.Ticker) ([]byte, error) {
	apiPath := fmt.Sprintf("api/fundamentals/%s?api_token=%s", ticker, c.apiKey)
	return c.process(ctx, apiPath, normalizeFundamentals, nil)
}

func (c *Client) HistoricalIntraDay(ctx context.Context, ticker config.Ticker, interval string, from time.Time, to time.Time) ([]byte, error) {
	apiPath := fmt.Sprintf("api/intraday/%s?api_token=%s&fmt=json&interval=%s&from=%d&to=%d", ticker, c.apiKey, interval, from.Unix(), to.Unix())

	// use a hook to insert the ticker within the dataset as well
	parsedHook := func(v any) error {
		m := v.(map[string]any)
		m["ticker"] = ticker.String()
		return nil
	}

	return c.process(ctx, apiPath, normalizeHistoricalIntraday, parsedHook)
}
