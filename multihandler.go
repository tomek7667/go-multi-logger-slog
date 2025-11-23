package logger

import (
	"context"
	"log/slog"
)

type MultiHandler struct {
	handlers []slog.Handler
}

func (mh *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (mh *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, handler := range mh.handlers {
		if err := handler.Handle(ctx, record); err != nil {
			return err
		}
	}
	return nil
}

func (mh *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, handler := range mh.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiHandler(newHandlers...)
}

func (mh *MultiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(mh.handlers))
	for i, handler := range mh.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return NewMultiHandler(newHandlers...)
}
