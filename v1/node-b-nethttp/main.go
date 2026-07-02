package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	port := getEnv("PORT", "8080")
	n, _ := strconv.Atoi(getEnv("SUBSCRIBER_COUNT", "100000"))

	log.Printf("seeding %d subscribers...", n)
	store := NewSubscriberStore(n)
	log.Printf("store ready, listening on :%s | GOMAXPROCS=%d", port, runtime.GOMAXPROCS(0))

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewHandler(store),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("shutting down...")
	_ = server.Close()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
