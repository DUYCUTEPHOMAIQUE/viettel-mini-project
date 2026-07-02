package main

import (
	"net/http"
	"strings"
)

var (
	bodyNotFound         = []byte(`{"error":"not_found"}`)
	bodyMethodNotAllowed = []byte(`{"error":"method_not_allowed"}`)
	bodyHealth           = []byte(`{"status":"ok"}`)
)

type Handler struct {
	store *SubscriberStore
}

func NewHandler(store *SubscriberStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	path := r.URL.Path
	switch {
	case strings.HasPrefix(path, "/subscriber/"):
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write(bodyMethodNotAllowed)
			return
		}
		supi := path[len("/subscriber/"):]
		h.getSubscriber(w, supi)

	case path == "/health":
		w.WriteHeader(http.StatusOK)
		w.Write(bodyHealth)

	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write(bodyNotFound)
	}
}

func (h *Handler) getSubscriber(w http.ResponseWriter, supi string) {
	data, ok := h.store.Get(supi)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write(bodyNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
