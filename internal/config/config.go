package config

import "time"

type Config struct {
	Endpoint string
	Insecure bool
	Timeout  time.Duration
	PAT      string
	OrgID    string
}

var Global Config
