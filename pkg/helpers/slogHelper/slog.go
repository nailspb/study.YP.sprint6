package slogHelper

import (
	"log/slog"
)

func GetErrAttr(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func ConfigureForMiddleware(log *slog.Logger, name string) *slog.Logger {
	log.Info("add middleware " + name)
	return log.With(slog.String("middleware", name))
}
