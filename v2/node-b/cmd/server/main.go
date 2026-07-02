package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"node-b/internal/config"
	"node-b/internal/logger"
	"node-b/internal/server"
	"node-b/internal/store"
)

func main() {
	configPath := flag.String("config", "assets/config.yaml", "path to config file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	log := logger.New()
	defer log.Sync()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("load config", zap.Error(err))
	}

	log.Info("seeding subscribers", zap.Int("count", cfg.SubscriberCount))
	subscribers := store.New(cfg.SubscriberCount)
	log.Info("store ready",
		zap.Int("port", cfg.Port),
		zap.Int("gomaxprocs", runtime.GOMAXPROCS(0)),
		zap.Bool("log_requests", cfg.LogRequests),
	)

	handler := server.NewHandler(subscribers, log, cfg.LogRequests)
	srv := &fasthttp.Server{
		Handler:            handler.HandleRequest,
		ReadTimeout:        cfg.ReadTimeout,
		WriteTimeout:       cfg.WriteTimeout,
		IdleTimeout:        cfg.IdleTimeout,
		DisableKeepalive:   false,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 30 * time.Second,
		Concurrency:        cfg.Concurrency,
	}

	go func() {
		if err := srv.ListenAndServe(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			log.Fatal("listen failed", zap.Error(err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Info("shutting down")
	_ = srv.Shutdown()
}
