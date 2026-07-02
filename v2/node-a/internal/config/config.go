package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host            string
	Port            int
	Concurrency     int
	Duration        time.Duration
	PoolSize        int
	SubscriberCount int
}

func (c Config) Target() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var raw struct {
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		Concurrency     int    `yaml:"concurrency"`
		Duration        string `yaml:"duration"`
		PoolSize        int    `yaml:"pool_size"`
		SubscriberCount int    `yaml:"subscriber_count"`
	}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return Config{}, err
	}

	duration, err := time.ParseDuration(raw.Duration)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Host:            raw.Host,
		Port:            raw.Port,
		Concurrency:     raw.Concurrency,
		Duration:        duration,
		PoolSize:        raw.PoolSize,
		SubscriberCount: raw.SubscriberCount,
	}, nil
}
