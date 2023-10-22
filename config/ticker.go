package config

import (
	"fmt"
	"strings"
)

type Ticker struct {
	Symbol string
	Market string
}

func (t Ticker) String() string {
	return t.Symbol + "." + t.Market
}

func (t Ticker) MarshalYAML() (any, error) {
	return t.String(), nil
}

func (t *Ticker) UnmarshalYAML(unmarshaller func(any) error) error {
	var v string
	if err := unmarshaller(&v); err != nil {
		return err
	}
	split := strings.Split(v, ".")

	if len(split) == 1 {
		*t = Ticker{
			Symbol: split[0],
			Market: "US",
		}
		return nil
	}

	if len(split) == 2 {
		*t = Ticker{
			Symbol: split[0],
			Market: split[1],
		}
		return nil
	}

	return fmt.Errorf("ticker is not in the <market>.<symbol> format")
}
