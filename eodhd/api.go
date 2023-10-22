package eodhd

import (
	"context"
	"fmt"
	"time"

	"github.com/bcap/stock-data/config"
	"github.com/bcap/stock-data/jq"
)

var normalizeExchanges = jq.LoadScript("normalizers/exchanges.jq")
var normalizeTickers = jq.LoadScript("normalizers/tickers.jq")
var normalizeFundamentals = jq.LoadScript("normalizers/fundamentals.jq")
var normalizeHistoricalIntraday = jq.LoadScript("normalizers/historical-intraday.jq")

func (c *Client) Exchanges(ctx context.Context, ts time.Time) ([]byte, error) {
	apiPath := fmt.Sprintf("api/exchanges-list/?api_token=%s&fmt=json", c.apiKey)
	return c.process(ctx, apiPath, normalizeExchanges, insertTs(ts))
}

func (c *Client) Tickers(ctx context.Context, exchange string, ts time.Time) ([]byte, error) {
	apiPath := fmt.Sprintf("api/exchange-symbol-list/%s?api_token=%s&fmt=json", exchange, c.apiKey)
	parsedHook := func(v any) error {
		m := v.(map[string]any)
		m["ExchangeGroup"] = exchange
		m["timestamp"] = ts.UTC().Unix()
		return nil
	}
	return c.process(ctx, apiPath, normalizeTickers, parsedHook)
}

func (c *Client) Fundamentals(ctx context.Context, ticker config.Ticker, ts time.Time) ([]byte, error) {
	apiPath := fmt.Sprintf("api/fundamentals/%s?api_token=%s", ticker, c.apiKey)
	return c.process(ctx, apiPath, normalizeFundamentals, insertTs(ts))
}

func (c *Client) HistoricalIntraDay(ctx context.Context, ticker config.Ticker, interval string, from time.Time, to time.Time) ([]byte, error) {
	apiPath := fmt.Sprintf("api/intraday/%s?api_token=%s&fmt=json&interval=%s&from=%d&to=%d", ticker, c.apiKey, interval, from.Unix(), to.Unix())
	parsedHook := func(v any) error {
		m := v.(map[string]any)
		m["ticker"] = ticker.String()
		return nil
	}
	return c.process(ctx, apiPath, normalizeHistoricalIntraday, parsedHook)
}

func insertTs(ts time.Time) jq.ParsedHook {
	return func(v any) error {
		m := v.(map[string]any)
		m["timestamp"] = ts.UTC().Unix()
		return nil
	}
}
