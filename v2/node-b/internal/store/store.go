package store

import "fmt"

type SubscriberStore struct {
	cache map[string][]byte
}

func New(n int) *SubscriberStore {
	cache := make(map[string][]byte, n)
	for i := 0; i < n; i++ {
		supi := fmt.Sprintf("imsi-00101%010d", i)
		cache[supi] = []byte(fmt.Sprintf(
			`{"supi":%q,"status":"ACTIVE","qos":"5G"}`, supi,
		))
	}
	return &SubscriberStore{cache: cache}
}

func (s *SubscriberStore) Get(supi string) ([]byte, bool) {
	b, ok := s.cache[supi]
	return b, ok
}
