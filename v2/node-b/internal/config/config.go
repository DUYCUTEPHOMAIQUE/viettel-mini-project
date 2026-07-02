package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port            int
	SubscriberCount int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	Concurrency     int
	LogRequests     bool
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var raw struct {
		Port            int    `yaml:"port"`
		SubscriberCount int    `yaml:"subscriber_count"`
		ReadTimeout     string `yaml:"read_timeout"`
		WriteTimeout    string `yaml:"write_timeout"`
		IdleTimeout     string `yaml:"idle_timeout"`
		Concurrency     int    `yaml:"concurrency"`
		LogRequests     bool   `yaml:"log_requests"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, err
	}

	readTimeout, err := time.ParseDuration(raw.ReadTimeout)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := time.ParseDuration(raw.WriteTimeout)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := time.ParseDuration(raw.IdleTimeout)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:            raw.Port,
		SubscriberCount: raw.SubscriberCount,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		IdleTimeout:     idleTimeout,
		Concurrency:     raw.Concurrency,
		LogRequests:     raw.LogRequests,
	}, nil
}
