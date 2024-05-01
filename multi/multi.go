package multi

import (
	"context"
	"log/slog"
)

// MultiHandler duplicates its Handle() to all handlers.
type MultiHandler struct {
	handlers []slog.Handler
}

func NewHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

// returns some handler is Enabled()
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, he := range h.handlers {
		if he != nil && he.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

// applies Handle() to all handlers.
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	var outerErr error
	for _, he := range h.handlers {
		err := he.Handle(ctx, r)
		if err != nil && outerErr == nil {
			outerErr = err
		}
	}

	if outerErr != nil {
		return outerErr
	}
	return nil
}

// applies WithAttrs() to all handlers.
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newhandlers := make([]slog.Handler, 0, len(h.handlers))

	for _, he := range h.handlers {
		newhandlers = append(newhandlers, he.WithAttrs(attrs))
	}

	return NewHandler(newhandlers...)
}

// applies WithGroup() to all handlers.
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	newhandlers := make([]slog.Handler, 0, len(h.handlers))

	for _, he := range h.handlers {
		newhandlers = append(newhandlers, he.WithGroup(name))
	}

	return NewHandler(newhandlers...)
}
