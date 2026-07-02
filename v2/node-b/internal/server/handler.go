package server

import (
	"bytes"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"node-b/internal/store"
)

var (
	prefixSubscriber     = []byte("/subscriber/")
	pathHealth           = []byte("/health")
	bodyNotFound         = []byte(`{"error":"not_found"}`)
	bodyMethodNotAllowed = []byte(`{"error":"method_not_allowed"}`)
	bodyHealth           = []byte(`{"status":"ok"}`)
	contentTypeJSON      = []byte("application/json; charset=utf-8")
)

type Handler struct {
	store       *store.SubscriberStore
	logger      *zap.Logger
	logRequests bool
}

func NewHandler(s *store.SubscriberStore, logger *zap.Logger, logRequests bool) *Handler {
	return &Handler{store: s, logger: logger, logRequests: logRequests}
}

func (h *Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
	if h.logRequests {
		start := time.Now()
		defer func() {
			h.logger.Info("request",
				zap.ByteString("method", ctx.Method()),
				zap.ByteString("path", ctx.Path()),
				zap.Int("status", ctx.Response.StatusCode()),
				zap.Duration("latency", time.Since(start)),
			)
		}()
	}

	path := ctx.Path()
	ctx.Response.Header.SetContentTypeBytes(contentTypeJSON)

	switch {
	case bytes.HasPrefix(path, prefixSubscriber):
		if !ctx.IsGet() {
			ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
			ctx.SetBody(bodyMethodNotAllowed)
			return
		}
		h.getSubscriber(ctx, string(path[len(prefixSubscriber):]))

	case bytes.Equal(path, pathHealth):
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(bodyHealth)

	default:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody(bodyNotFound)
	}
}

func (h *Handler) getSubscriber(ctx *fasthttp.RequestCtx, supi string) {
	data, ok := h.store.Get(supi)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.SetBody(bodyNotFound)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(data)
}
