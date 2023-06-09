package config

import (
	"errors"
	"fmt"
	"os"

	gostrings "github.com/555f/go-strings"
)

func New() (c *Config, errs []error) {
	c = &Config{}
	if s, ok := os.LookupEnv("ADDR"); ok {
		c.Addr = s
		if c.Addr == "" {
			errs = append(errs, errors.New("env ADDR empty"))
		}
	} else {
		errs = append(errs, errors.New("env ADDR not set"))
	}
	if s, ok := os.LookupEnv("PORT"); ok {
		v, err := gostrings.ParseInt[int](s, 10, 64)
		if err != nil {
			errs = append(errs, fmt.Errorf("env PORT failed parse: %w", err))
		}
		c.Port = v
		if c.Port == 0 {
			errs = append(errs, errors.New("env PORT empty"))
		}
	} else {
		errs = append(errs, errors.New("env PORT not set"))
	}
	return
}
