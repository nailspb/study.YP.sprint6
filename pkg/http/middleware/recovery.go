package middleware

import (
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/helpers/slogHelper"
	"github.com/Yandex-Practicum/go-rest-api-homework/pkg/http/routing"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(log *slog.Logger) routing.MiddlewareFunc {
	log = slogHelper.ConfigureForMiddleware(log, "Recovery")
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(http.StatusText(http.StatusInternalServerError), slogHelper.GetErrAttr(err.(error)), slog.String("stack", string(debug.Stack())))
					http.Error(w, http.StatusText(http.StatusInternalServerError)+"\r\n\r\n"+string(debug.Stack()), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, req)
		}
	}
}
