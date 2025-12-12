package config

import "github.com/caarlos0/env/v11"

func parseEnv[T any](opts env.Options) (T, error) {
	var t T

	if err := env.Parse(&t); err != nil {
		return t, err
	}

	env.ParseWithOptions(&t, opts)

	return t, nil
}
