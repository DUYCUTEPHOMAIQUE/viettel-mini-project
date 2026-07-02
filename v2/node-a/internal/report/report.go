package report

import (
	"fmt"

	"node-a/internal/client"
	"node-a/internal/loadtest"
)

func Print(result loadtest.Result, stats client.Stats) {
	fmt.Println("------------------------------------------")
	fmt.Printf("requests      : %d\n", result.Requests)
	fmt.Printf("errors        : %d (%.3f%%)\n", result.Errors, result.ErrorRate())
	fmt.Printf("throughput    : %.0f rps\n", result.Throughput())
	fmt.Printf("p50 latency   : %s\n", result.Percentile(0.50))
	fmt.Printf("p95 latency   : %s\n", result.Percentile(0.95))
	fmt.Printf("p99 latency   : %s\n", result.Percentile(0.99))
	fmt.Printf("new tcp conns : %d\n", stats.NewConns)
	fmt.Printf("reused conns  : %d\n", stats.ReusedConns)
	fmt.Println("------------------------------------------")
}
