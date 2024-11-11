// Code generated by GG version . DO NOT EDIT.

//go:build !gg

package config

import (
	"fmt"
	env "github.com/555f/gg/examples/env"
	gostrings "github.com/555f/go-strings"
	gomultierror "github.com/hashicorp/go-multierror"
	"os"
)

func New() (c *env.Config, errs error) {
	c = &env.Config{}
	if s, ok := os.LookupEnv("PORT"); ok {

		v, err := gostrings.ParseInt[int](s, 10, 64)
		if err != nil {
			errs = gomultierror.Append(errs, fmt.Errorf("env PORT failed parse: %w", err))
		}
		c.Port = v
	}
	return
}
