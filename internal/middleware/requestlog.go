package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// RequestLog returns a middleware that logs each incoming requests using the
// given logger.
func RequestLog(lg *logrus.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return &handler{delegate: h, logger: lg}
	}
}

type handler struct {
	delegate http.Handler
	logger   *logrus.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	severity := logrus.InfoLevel
	if r.Method == http.MethodGet && (r.URL.Path == "/health" || r.URL.Path == "/version") {
		severity = logrus.DebugLevel
	}

	w2 := &responseWriter{ResponseWriter: w}
	h.delegate.ServeHTTP(w2, r)
	s := severity
	if w2.status >= 500 {
		s = logrus.ErrorLevel
	}

	h.logger.WithFields(logrus.Fields{
		"method":       r.Method,
		"url":          r.URL,
		"responseSize": w2.length,
		"userAgent":    r.UserAgent(),
		"remoteIP":     remoteAddr(r),
		"referer":      r.Referer(),
		"status":       translateStatus(w2.status),
		"latency":      time.Since(start),
	}).Log(s)
}

type responseWriter struct {
	http.ResponseWriter

	status int
	length int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ResponseWriter.Write(p)
	rw.length += n
	return
}

func translateStatus(code int) int {
	if code == 0 {
		return http.StatusOK
	}
	return code
}

func remoteAddr(r *http.Request) string {
	var addrs []string

	if r.Header.Get("x-forwarded-for") != "" {
		addrs = append(addrs, r.Header.Get("x-forwarded-for"))
	}

	addrs = append(addrs, r.RemoteAddr)

	return strings.Join(addrs, ",")
}
