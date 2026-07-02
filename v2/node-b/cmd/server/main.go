package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"

	"node-b/internal/config"
	"node-b/internal/server"
	"node-b/internal/store"
)

func main() {
	configPath := flag.String("config", "assets/config.yaml", "path to config file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	log.Printf("seeding %d subscribers...", cfg.SubscriberCount)
	subscribers := store.New(cfg.SubscriberCount)
	log.Printf("store ready, listening on :%d | GOMAXPROCS=%d", cfg.Port, runtime.GOMAXPROCS(0))

	handler := server.NewHandler(subscribers)
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
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutting down...")
	_ = srv.Shutdown()
}
