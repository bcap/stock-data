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
	ApiKey               string `yaml:"api-key"`
	MaxRequestsPerMinute int    `yaml:"max-requests-per-min"`
}

type AWS struct {
	Region  string `yaml:"region"`
	Profile string `yaml:"profile"`
}

func (c *Config) UnmarshalYAML(unmarshaller func(any) error) error {
	type raw Config
	v := raw{}
	if err := unmarshaller(&v); err != nil {
		return err
	}
	*c = Config(v)
	c.Tickers = unique(c.Tickers, func(v *Ticker) string { return v.String() })
	c.Exchanges = unique(c.Exchanges, func(v *string) string { return *v })
	return nil
}

func unique[T any](input []T, keyFn func(v *T) string) []T {
	if input == nil {
		return nil
	}
	seen := map[string]struct{}{}
	result := []T{}
	for _, v := range input {
		if _, ok := seen[keyFn(&v)]; ok {
			continue
		}
		result = append(result, v)
	}
	return result
}
