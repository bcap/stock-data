package config

import (
	"fmt"
	"time"
)

type empty struct{}

type LoadOption string

type LoadOptions map[LoadOption]empty

const (
	LoadExchanges          LoadOption = "exchanges"
	LoadFundamentals       LoadOption = "fundamentals"
	LoadHistoricalIntraday LoadOption = "historical-intraday"
)

var allLoadOptions LoadOptions = LoadOptions{
	LoadExchanges:          empty{},
	LoadFundamentals:       empty{},
	LoadHistoricalIntraday: empty{},
}

var Default = Config{
	Tickers: []Ticker{},
	Load: LoadOptions{
		LoadFundamentals:       {},
		LoadHistoricalIntraday: {},
	},
	EODHD:       EODHD{ApiKey: "demo"},
	AWS:         AWS{Region: "us-east-1", Profile: "default"},
	MaxParallel: 10,
}

type Config struct {
	Tickers            []Ticker           `yaml:"tickers"`
	Load               LoadOptions        `yaml:"load"`
	Exchanges          Exchanges          `yaml:"exchanges"`
	Fundamentals       Fundamentals       `yaml:"fundamentals"`
	HistoricalIntraday HistoricalIntraday `yaml:"historical-intraday"`
	EODHD              EODHD              `yaml:"eodhd"`
	AWS                AWS                `yaml:"aws"`
	MaxParallel        int64              `yaml:"max-parallelism"`
}

type S3Stored struct {
	S3Bucket Path `yaml:"s3-bucket"`
	S3Prefix Path `yaml:"s3-prefix"`
}

type Exchanges struct {
	S3Stored `yaml:",inline"`
}

type Fundamentals struct {
	S3Stored `yaml:",inline"`
}

type HistoricalIntraday struct {
	S3Stored  `yaml:",inline"`
	TimeRange TimeRange `yaml:"time-range"`
	Interval  string    `yaml:"interval"`
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

func (l *LoadOptions) UnmarshalYAML(unmarshaller func(any) error) error {
	var values []string
	if err := unmarshaller(&values); err != nil {
		return err
	}
	result := LoadOptions{}
	for _, value := range values {
		lo := LoadOption(value)
		if _, ok := allLoadOptions[lo]; !ok {
			return fmt.Errorf("unknown load option %s", value)
		}
		result[LoadOption(value)] = empty{}
	}
	*l = result
	return nil
}
