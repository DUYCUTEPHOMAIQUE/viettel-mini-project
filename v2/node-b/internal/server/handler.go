package server

import (
	"bytes"

	"github.com/valyala/fasthttp"

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
	store *store.SubscriberStore
}

func NewHandler(s *store.SubscriberStore) *Handler {
	return &Handler{store: s}
}

func (h *Handler) HandleRequest(ctx *fasthttp.RequestCtx) {
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
