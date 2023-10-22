package config

import (
	"time"
)

type Config struct {
	Exchanges          []string           `yaml:"exchanges"`
	Tickers            []Ticker           `yaml:"tickers"`
	ListExchanges      ListExchanges      `yaml:"list-exchanges"`
	ListTickers        ListTickers        `yaml:"list-tickers"`
	Fundamentals       Fundamentals       `yaml:"fundamentals"`
	HistoricalIntraday HistoricalIntraday `yaml:"historical-intraday"`
	EODHD              EODHD              `yaml:"eodhd"`
	AWS                AWS                `yaml:"aws"`
	MaxParallel        int64              `yaml:"max-parallelism"`
}

type FetchConfig struct {
	Enabled  bool `yaml:"enabled"`
	S3Bucket Path `yaml:"s3-bucket"`
	S3Prefix Path `yaml:"s3-prefix"`
}

type ListExchanges struct {
	FetchConfig `yaml:",inline"`
}

type ListTickers struct {
	FetchConfig `yaml:",inline"`
}

type Fundamentals struct {
	FetchConfig `yaml:",inline"`
}

type HistoricalIntraday struct {
	FetchConfig `yaml:",inline"`
	TimeRange   TimeRange `yaml:"time-range"`
	Interval    string    `yaml:"interval"`
}

type TimeRange struct {
	From  Time          `yaml:"from"`
	To    Time          `yaml:"to"`
	Split time.Duration `yaml:"split"`
}

type EODHD struct {
	ApiKey string `yaml:"api-key"`
}

type AWS struct {
	Region  string `yaml:"region"`
	Profile string `yaml:"profile"`
}
