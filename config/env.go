package config

import (
	"strings"
)

const (
	Local string = "LOCAL"
	Dev   string = "DEV"
	UAT   string = "UAT"
	Prod  string = "PROD"
)

func (c *Config) IsLocalEnv() bool {
	return strings.ToUpper(c.Env) == Local
}

func (c *Config) IsDevEnv() bool {
	return strings.ToUpper(c.Env) == Dev
}

func (c *Config) IsUATEnv() bool {
	return strings.ToUpper(c.Env) == UAT
}

func (c *Config) IsProdEnv() bool {
	return strings.ToUpper(c.Env) == Prod
}
