package leveled

import (
	"context"
	"log/slog"
	"os"
)

// LeveledHandler has handlers for levels.
type LeveledHandler struct {
	defaultHandler, debugHandler, errorHandler, infoHandler, warnHandler slog.Handler
}

type LeveledOption func(*LeveledHandler)

func Debug(h slog.Handler) func(*LeveledHandler) {
	return func(lh *LeveledHandler) {
		lh.debugHandler = h
	}
}

func Error(h slog.Handler) func(*LeveledHandler) {
	return func(lh *LeveledHandler) {
		lh.errorHandler = h
	}
}

func Info(h slog.Handler) func(*LeveledHandler) {
	return func(lh *LeveledHandler) {
		lh.infoHandler = h
	}
}

func Warn(h slog.Handler) func(*LeveledHandler) {
	return func(lh *LeveledHandler) {
		lh.warnHandler = h
	}
}

// If defaultHandler == nil then slog.NewTextHandler(os.Stdout, nil) is used as a default handler.
//
// Use Debug(), Error(), Info(), Warn() to handle each level.
func NewHandler(defaultHandler slog.Handler, lopts ...LeveledOption) *LeveledHandler {
	h := LeveledHandler{
		defaultHandler: defaultHandler,
	}
	for _, o := range lopts {
		o(&h)
	}

	if h.defaultHandler == nil {
		h.defaultHandler = slog.NewTextHandler(os.Stdout, nil)
	}

	return &h
}

// returns some handler is Enabled()
func (h *LeveledHandler) Enabled(ctx context.Context, level slog.Level) bool {
	hh := []slog.Handler{
		h.defaultHandler,
		h.debugHandler,
		h.errorHandler,
		h.infoHandler,
		h.warnHandler,
	}
	for _, he := range hh {
		if he != nil && he.Enabled(ctx, level) {
			return true
		}
	}

	return false
}

// applies Handle() to appropriate level handler.
func (h *LeveledHandler) Handle(ctx context.Context, r slog.Record) error {
	hh := h.defaultHandler
	if r.Level == slog.LevelDebug && h.debugHandler != nil {
		hh = h.debugHandler
	}
	if r.Level == slog.LevelError && h.errorHandler != nil {
		hh = h.errorHandler
	}
	if r.Level == slog.LevelInfo && h.infoHandler != nil {
		hh = h.infoHandler
	}
	if r.Level == slog.LevelWarn && h.warnHandler != nil {
		hh = h.warnHandler
	}

	return hh.Handle(ctx, r)
}

// applies WithAttrs() to all handlers.
func (h *LeveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newh := NewHandler(h.defaultHandler.WithAttrs(attrs))

	if h.debugHandler != nil {
		newh.debugHandler = h.debugHandler.WithAttrs(attrs)
	}
	if h.errorHandler != nil {
		newh.errorHandler = h.errorHandler.WithAttrs(attrs)
	}
	if h.infoHandler != nil {
		newh.infoHandler = h.infoHandler.WithAttrs(attrs)
	}
	if h.warnHandler != nil {
		newh.warnHandler = h.warnHandler.WithAttrs(attrs)
	}

	return newh
}

// applies WithGroup() to all handlers.
func (h *LeveledHandler) WithGroup(name string) slog.Handler {
	newh := NewHandler(h.defaultHandler.WithGroup(name))

	if h.debugHandler != nil {
		newh.debugHandler = h.debugHandler.WithGroup(name)
	}
	if h.errorHandler != nil {
		newh.errorHandler = h.errorHandler.WithGroup(name)
	}
	if h.infoHandler != nil {
		newh.infoHandler = h.infoHandler.WithGroup(name)
	}
	if h.warnHandler != nil {
		newh.warnHandler = h.warnHandler.WithGroup(name)
	}

	return newh
}
