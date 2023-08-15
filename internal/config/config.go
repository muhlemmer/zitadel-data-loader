package config

import "time"

type Config struct {
	Endpoint      string
	Insecure      bool
	Timeout       time.Duration
	APIUserTokens map[string]string
}

var Global Config
