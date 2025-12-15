package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env      string `env:"ENV" default:"LOCAL"`
	Server   Server
	Database Database
	Service  Service
}

type Server struct {
	Hostname string `env:"HOSTNAME"`
	Port     string `env:"PORT"`
}

type Database struct {
	PostgresURL string `env:"POSTGRES_URL"`
}

type Service struct {
	ReferralService string `env:"REFERRAL_SERVICE"`
}

var once sync.Once
var config Config

func prefix(e string) string {
	if e == "" {
		return ""
	}

	return fmt.Sprintf("%s_", e)
}

func C(envPrefix string) Config {
	once.Do(func() {
		opts := env.Options{
			Prefix: prefix(envPrefix),
		}

		var err error
		config, err = parseEnv[Config](opts)
		if err != nil {
			log.Fatal(err)
		}
	})

	return config
}
