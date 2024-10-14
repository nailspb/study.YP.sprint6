package middleware

import (
	"fmt"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/helpers/slogHelper"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/http/routing"
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(log *slog.Logger) routing.MiddlewareFunc {
	log = slogHelper.ConfigureForMiddleware(log, "Logging")
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			start := time.Now()
			ip := req.Header.Get("X-Real-IP")
			if ip == "" {
				ip = req.RemoteAddr
			}
			log.Info("Request start",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("query", req.URL.RawQuery),
				slog.String("ip", ip),
				slog.String("user-agent", req.UserAgent()),
			)
			recorder := &statusRecorder{
				ResponseWriter: w,
				Status:         200,
			}
			next.ServeHTTP(recorder, req)
			duration := time.Since(start)
			log.Info("Request end",
				slog.Int("StatusCode", recorder.Status),
				slog.String("duration", fmt.Sprintf("%d us", duration.Microseconds())),
			)

		}
	}
}
