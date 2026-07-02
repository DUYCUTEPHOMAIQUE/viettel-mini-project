package loadtest

import (
	"sort"
	"time"
)

type Result struct {
	Requests  int64
	Errors    int64
	Elapsed   time.Duration
	Latencies []time.Duration
}

func mergeResults(partials []Result, elapsed time.Duration) Result {
	total := Result{Elapsed: elapsed}
	for _, p := range partials {
		total.Requests += p.Requests
		total.Errors += p.Errors
		total.Latencies = append(total.Latencies, p.Latencies...)
	}
	sort.Slice(total.Latencies, func(i, j int) bool {
		return total.Latencies[i] < total.Latencies[j]
	})
	return total
}

func (r Result) Throughput() float64 {
	if r.Elapsed <= 0 {
		return 0
	}
	return float64(r.Requests) / r.Elapsed.Seconds()
}

func (r Result) ErrorRate() float64 {
	total := r.Requests + r.Errors
	if total == 0 {
		return 0
	}
	return float64(r.Errors) / float64(total) * 100
}

func (r Result) Percentile(p float64) time.Duration {
	if len(r.Latencies) == 0 {
		return 0
	}
	idx := int(float64(len(r.Latencies)-1) * p)
	return r.Latencies[idx]
}
