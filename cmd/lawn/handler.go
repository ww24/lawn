package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ww24/lawn"
)

const (
	negativeCacheDuration = 30 * time.Second
)

type handler struct {
	cli    *lawn.Client
	maxAge time.Duration
}

func newHandler(cli *lawn.Client, maxAge time.Duration) http.Handler {
	return &handler{
		cli:    cli,
		maxAge: maxAge,
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.Handle("/", h.handler())

	// middleware
	handler := h.accesslog(mux)
	handler.ServeHTTP(w, r)
}

func (h *handler) accesslog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("request mathod:%s, path:%s, userAgent:%s\n", r.Method, r.URL.EscapedPath(), r.UserAgent())
		next.ServeHTTP(w, r)
	})
}

func (h *handler) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Content-Type", "image/svg+xml")
		header.Set("Cache-Control", fmt.Sprintf("public, max-age=0, s-maxage=%d", int(h.maxAge.Seconds())))

		switch r.Method {
		case http.MethodHead:
			return
		case http.MethodOptions:
			header.Set("Allow", "OPTIONS, HEAD, GET")
			return
		case http.MethodGet:
		default:
			h.sendError(w, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		if err := h.cli.Fetch(ctx, w, username); err != nil {
			h.sendError(w, http.StatusInternalServerError)
			return
		}
	})
}

func (h *handler) sendError(w http.ResponseWriter, status int) {
	header := w.Header()
	header.Set("Content-Type", "text/plain")
	// negative cache
	header.Set("Cache-Control", fmt.Sprintf("public, max-age=0, s-maxage=%d", int(negativeCacheDuration.Seconds())))
	w.WriteHeader(status)
}
