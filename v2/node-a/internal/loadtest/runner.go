package loadtest

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Fetcher interface {
	GetSubscriber(ctx context.Context, supi string) (int, error)
}

type Runner struct {
	fetcher         Fetcher
	workers         int
	duration        time.Duration
	subscriberCount int
}

func NewRunner(fetcher Fetcher, workers int, duration time.Duration, subscriberCount int) *Runner {
	return &Runner{
		fetcher:         fetcher,
		workers:         workers,
		duration:        duration,
		subscriberCount: subscriberCount,
	}
}

func (r *Runner) Run() Result {
	ctx, cancel := context.WithTimeout(context.Background(), r.duration)
	defer cancel()

	perWorker := make([]Result, r.workers)
	var wg sync.WaitGroup
	wg.Add(r.workers)

	start := time.Now()
	for id := 0; id < r.workers; id++ {
		go func(id int) {
			defer wg.Done()
			perWorker[id] = r.runWorker(ctx, id)
		}(id)
	}
	wg.Wait()
	elapsed := time.Since(start)

	return mergeResults(perWorker, elapsed)
}

func (r *Runner) runWorker(ctx context.Context, id int) Result {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
	stats := Result{Latencies: make([]time.Duration, 0, 1024)}

	for ctx.Err() == nil {
		supi := makeSUPI(rng.Intn(r.subscriberCount))

		start := time.Now()
		status, err := r.fetcher.GetSubscriber(ctx, supi)
		elapsed := time.Since(start)

		if ctx.Err() != nil {
			break
		}
		if err != nil || status != 200 {
			stats.Errors++
			continue
		}
		stats.Requests++
		stats.Latencies = append(stats.Latencies, elapsed)
	}
	return stats
}

func makeSUPI(i int) string {
	return fmt.Sprintf("imsi-00101%010d", i)
}
