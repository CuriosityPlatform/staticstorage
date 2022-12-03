package config

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

type Parser struct{}

func (p Parser) ParseConfig(reader io.Reader) (Config, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return Config{}, errors.Wrapf(err, "failed to read config")
	}

	c := Config{}

	err = json.Unmarshal(bytes, &c)
	return c, errors.Wrapf(err, "failed to parse config")
}
