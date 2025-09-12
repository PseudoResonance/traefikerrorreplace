// Package cloudflarewarp Traefik Plugin.
package cloudflarewarp

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
	MatchStatus   []int `json:"matchStatus,omitempty"`
	ReplaceStatus int   `json:"replaceStatus,omitempty"`
	Debug         bool  `json:"debug,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		MatchStatus:   []int{},
		ReplaceStatus: 200, //nolint:mnd
		Debug:         false,
	}
}

// StatusCodeReplacer is a plugin that replaces status codes.
type StatusCodeReplacer struct {
	next          http.Handler
	name          string
	MatchStatus   []int
	ReplaceStatus int
	Debug         bool
}

// New created a new plugin.
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	statusCodeReplacer := &StatusCodeReplacer{
		next: next,
		name: name,
	}

	if config.MatchStatus != nil {
		statusCodeReplacer.MatchStatus = append(statusCodeReplacer.MatchStatus, config.MatchStatus...)
	}

	statusCodeReplacer.ReplaceStatus = config.ReplaceStatus

	statusCodeReplacer.Debug = config.Debug

	return statusCodeReplacer, nil
}

/**
 * The following section is licensed under the following
 * Apache License, Version 2.0, January 2004
 * Copyright 2020 Containous SAS
 * Copyright 2020 Traefik Labs
 *
 * Inspired by the following plugin
 * https://github.com/XciD/traefik-plugin-rewrite-headers/blob/master/rewrite_headers.go
 */

type responseWriter struct {
	writer        http.ResponseWriter
	matchStatus   []int
	replaceStatus int
}

func (r *StatusCodeReplacer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wrappedWriter := &responseWriter{
		writer:        rw,
		matchStatus:   r.MatchStatus,
		replaceStatus: r.ReplaceStatus,
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

/**
 * End section
 */
