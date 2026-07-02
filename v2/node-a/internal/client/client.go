package client

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"sync/atomic"
	"time"
)

type Stats struct {
	NewConns    int64
	ReusedConns int64
}

type Client struct {
	httpClient  *http.Client
	baseURL     string
	newConns    atomic.Int64
	reusedConns atomic.Int64
}

func New(baseURL string, poolSize int) *Client {
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	transport := &http.Transport{
		DialContext:         dialer.DialContext,
		MaxIdleConns:        poolSize,
		MaxIdleConnsPerHost: poolSize,
		MaxConnsPerHost:     poolSize,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
		ForceAttemptHTTP2:   false,
	}
	return &Client{
		httpClient: &http.Client{Transport: transport, Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

func (c *Client) GetSubscriber(ctx context.Context, supi string) (int, error) {
	url := c.baseURL + "/subscriber/" + supi
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			if info.Reused {
				c.reusedConns.Add(1)
			} else {
				c.newConns.Add(1)
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	return resp.StatusCode, nil
}

func (c *Client) Stats() Stats {
	return Stats{
		NewConns:    c.newConns.Load(),
		ReusedConns: c.reusedConns.Load(),
	}
}
