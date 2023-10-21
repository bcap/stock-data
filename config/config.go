package config

import (
	"fmt"
	"time"
)

type LoadOption string

type LoadOptions map[LoadOption]struct{}

const (
	LoadFundamentals       LoadOption = "fundamentals"
	LoadHistoricalIntraday LoadOption = "historical-intraday"
)

var allLoadOptions LoadOptions = LoadOptions{
	LoadFundamentals:       struct{}{},
	LoadHistoricalIntraday: struct{}{},
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
	Fundamentals       Fundamentals       `yaml:"fundamentals"`
	HistoricalIntraday HistoricalIntraday `yaml:"historical-intraday"`
	EODHD              EODHD              `yaml:"eodhd"`
	AWS                AWS                `yaml:"aws"`
	MaxParallel        int64              `yaml:"max-parallelism"`
}

type Fundamentals struct {
	S3Bucket Path `yaml:"s3-bucket"`
	S3Prefix Path `yaml:"s3-prefix"`
}

type HistoricalIntraday struct {
	TimeRange TimeRange `yaml:"time-range"`
	Interval  string    `yaml:"interval"`
	S3Bucket  Path      `yaml:"s3-bucket"`
	S3Prefix  Path      `yaml:"s3-prefix"`
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
		result[LoadOption(value)] = struct{}{}
	}
	*l = result
	return nil
}
