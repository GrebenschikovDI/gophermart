package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

type loggerKey struct{}

func LoggerMiddleware(log *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := context.WithValue(r.Context(), loggerKey{}, log)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func LoggerFromContext(ctx context.Context) *logrus.Logger {
	return ctx.Value(loggerKey{}).(*logrus.Logger)
}

func LogError(w http.ResponseWriter, r *http.Request, err error) {
	LoggerFromContext(r.Context()).WithFields(logrus.Fields{
		"error":  err.Error(),
		"method": r.Method,
		"path":   r.URL.Path,
	}).Error("Error in handler")
}
