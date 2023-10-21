package config

import (
	"io"
	"os"

	"github.com/barkimedes/go-deepcopy"
	"gopkg.in/yaml.v2"
)

func Load(path string) (Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	cfg := deepcopy.MustAnything(Default).(Config)
	if err := Parse(f, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Parse(reader io.Reader, config *Config) error {
	return yaml.NewDecoder(reader).Decode(config)
}
