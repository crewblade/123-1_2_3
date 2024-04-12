package slogempty

import (
	"context"
	"log/slog"
)

type EmptyHandler struct{}

func NewEmptyLogger() *slog.Logger {
	return slog.New(NewEmptyHandler())
}

func NewEmptyHandler() *EmptyHandler {
	return &EmptyHandler{}
}

func (h *EmptyHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}

func (h *EmptyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *EmptyHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *EmptyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
