package config

import (
	"fmt"
	"strings"
)

type Ticker struct {
	Market string
	Symbol string
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
	if len(split) != 2 {
		return fmt.Errorf("ticker is not in the <market>.<symbol> format")
	}
	*t = Ticker{
		Market: split[0],
		Symbol: split[1],
	}
	return nil
}
