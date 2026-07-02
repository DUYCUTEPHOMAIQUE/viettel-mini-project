package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/valyala/fasthttp"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := getEnv("PORT", "8080")
	n, _ := strconv.Atoi(getEnv("SUBSCRIBER_COUNT", "100000"))

	log.Printf("seeding %d subscribers...", n)
	store := NewSubscriberStore(n)
	log.Printf("store ready, listening on :%s | GOMAXPROCS=%d", port, runtime.GOMAXPROCS(0))

	server := &fasthttp.Server{
		Handler:            NewHandler(store).HandleRequest,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       10 * time.Second,
		IdleTimeout:        60 * time.Second,
		DisableKeepalive:   false,
		TCPKeepalive:       true,
		TCPKeepalivePeriod: 30 * time.Second,
		Concurrency:        256 * 1024,
	}

	go func() {
		if err := server.ListenAndServe(":" + port); err != nil {
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutting down...")
	_ = server.Shutdown()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
