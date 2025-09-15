// Package traefikerrorreplace Traefik Plugin.
package traefikerrorreplace

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"slices"
)

// Config the plugin configuration.
type Config struct {
	Debug         bool  `json:"debug,omitempty"`
	MatchStatus   []int `json:"matchStatus,omitempty"`
	ReplaceStatus int   `json:"replaceStatus,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Debug:         false,
		MatchStatus:   []int{},
		ReplaceStatus: 200, //nolint:mnd
	}
}

// StatusCodeReplacer is a plugin that replaces status codes.
type StatusCodeReplacer struct {
	next          http.Handler
	name          string
	debug         bool
	matchStatus   []int
	replaceStatus int
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	statusCodeReplacer := &StatusCodeReplacer{
		next:          next,
		name:          name,
		debug:         config.Debug,
		matchStatus:   config.MatchStatus,
		replaceStatus: config.ReplaceStatus,
	}

	if statusCodeReplacer.debug {
		fmt.Printf("Debug printing enabled!")
	}

	return statusCodeReplacer, nil
}

type responseWriter struct {
	writer        http.ResponseWriter
	matchStatus   []int
	replaceStatus int
}

func (r *StatusCodeReplacer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:        rw,
		matchStatus:   r.matchStatus,
		replaceStatus: r.replaceStatus,
	}

	r.next.ServeHTTP(wrappedWriter, req)
}

func (r *responseWriter) Header() http.Header {
	return r.writer.Header()
}

func (r *responseWriter) Write(bytes []byte) (int, error) {
	return r.writer.Write(bytes)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	if slices.Contains(r.matchStatus, statusCode) {
		r.writer.WriteHeader(r.replaceStatus)
	}

	r.writer.WriteHeader(statusCode)
}

func (r *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.writer.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.writer)
	}

	return hijacker.Hijack()
}

func (r *responseWriter) Flush() {
	if flusher, ok := r.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}
