package handler

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	lr "github.com/Wallruzz9114/bookey/util/logger"
)

// Handler ...
type Handler struct {
	handler http.Handler
	logger  *lr.Logger
}

// NewHandler ...
func NewHandler(h http.HandlerFunc, l *lr.Logger) *Handler {
	return &Handler{handler: h, logger: l}
}

// ServeHTTP ...
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logEntry := &logEntry{
		ReceivedTime:      startTime,
		RequestMethod:     r.Method,
		RequestURL:        r.URL.String(),
		RequestHeaderSize: headerSize(r.Header),
		UserAgent:         r.UserAgent(),
		Referer:           r.Referer(),
		Proto:             r.Proto,
		RemoteIP:          ipFromHostPort(r.RemoteAddr),
	}

	if address, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		logEntry.ServerIP = ipFromHostPort(address.String())
	}

	r2 := new(http.Request)
	*r2 = *r
	rcc := &readCounterCloser{r: r.Body}
	r2.Body = rcc
	w2 := &responseStats{w: w}

	h.handler.ServeHTTP(w2, r2)

	logEntry.Latency = time.Since(startTime)

	if rcc.err == nil && rcc.r != nil {
		// If the handler hasn't encountered an error in the Body (like EOF),
		// then consume the rest of the Body to provide an accurate rcc.n.
		io.Copy(ioutil.Discard, rcc)
	}

	logEntry.RequestBodySize = rcc.n
	logEntry.Status = w2.code

	if logEntry.Status == 0 {
		logEntry.Status = http.StatusOK
	}

	logEntry.ResponseHeaderSize, logEntry.ResponseBodySize = w2.size()

	h.logger.Info().
		Time("received_time", logEntry.ReceivedTime).
		Str("method", logEntry.RequestMethod).
		Str("url", logEntry.RequestURL).
		Int64("header_size", logEntry.RequestHeaderSize).
		Int64("body_size", logEntry.RequestBodySize).
		Str("agent", logEntry.UserAgent).
		Str("referer", logEntry.Referer).
		Str("proto", logEntry.Proto).
		Str("remote_ip", logEntry.RemoteIP).
		Str("server_ip", logEntry.ServerIP).
		Int("status", logEntry.Status).
		Int64("resp_header_size", logEntry.ResponseHeaderSize).
		Int64("resp_body_size", logEntry.ResponseBodySize).
		Dur("latency", logEntry.Latency).
		Msg("")
}
