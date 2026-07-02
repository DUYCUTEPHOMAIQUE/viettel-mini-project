package main

import (
	"flag"
	"fmt"
	"log"

	"node-a/internal/client"
	"node-a/internal/config"
	"node-a/internal/loadtest"
	"node-a/internal/report"
)

func main() {
	configPath := flag.String("config", "assets/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	c := client.New(cfg.Target(), cfg.PoolSize)
	runner := loadtest.NewRunner(c, cfg.Concurrency, cfg.Duration, cfg.SubscriberCount)

	fmt.Printf("target=%s concurrency=%d duration=%s pool=%d\n",
		cfg.Target(), cfg.Concurrency, cfg.Duration, cfg.PoolSize)

	result := runner.Run()
	report.Print(result, c.Stats())
}
