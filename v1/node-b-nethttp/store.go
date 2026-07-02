package main

import "fmt"

// SubscriberStore pre-serializes JSON at startup.
// Read-only after init → safe concurrent access, zero alloc on hot path.
type SubscriberStore struct {
	cache map[string][]byte
}

func NewSubscriberStore(n int) *SubscriberStore {
	cache := make(map[string][]byte, n)
	for i := 0; i < n; i++ {
		supi := fmt.Sprintf("imsi-00101%010d", i)
		cache[supi] = []byte(fmt.Sprintf(
			`{"supi":%q,"status":"REGISTERED","plmnId":"00101"}`, supi,
		))
	}
	return &SubscriberStore{cache: cache}
}

func (s *SubscriberStore) Get(supi string) ([]byte, bool) {
	b, ok := s.cache[supi]
	return b, ok
}
