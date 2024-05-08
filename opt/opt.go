package opt

import (
	"context"
	"io"
	"log/slog"
	"sync"
)

type NewHandlerFunc func(*slog.HandlerOptions) slog.Handler

func textHandler(w io.Writer) NewHandlerFunc {
	return func(opts *slog.HandlerOptions) slog.Handler {
		return slog.NewTextHandler(w, opts)
	}
}

// a wrapper of slog.NewTextHandler
func NewTextHandler(w io.Writer, opts *slog.HandlerOptions) *OptHandler {
	return NewHandler(textHandler(w), opts)
}

func jsonHandler(w io.Writer) NewHandlerFunc {
	return func(opts *slog.HandlerOptions) slog.Handler {
		return slog.NewJSONHandler(w, opts)
	}
}

// a wrapper of slog.NewJSONHandler
func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) *OptHandler {
	return NewHandler(jsonHandler(w), opts)
}

// OptHandler is a wrapper handler that has some change-options methods.
type OptHandler struct {
	inner slog.Handler

	newfunc NewHandlerFunc

	mu   *sync.Mutex
	opts slog.HandlerOptions
}

func NewHandler(newfunc NewHandlerFunc, opts *slog.HandlerOptions) *OptHandler {
	h := &OptHandler{
		inner:   newfunc(opts),
		newfunc: newfunc,
		mu:      &sync.Mutex{},
	}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}

	return h
}

func (h *OptHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *OptHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.inner.Handle(ctx, r)
}

func (h *OptHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	return &OptHandler{
		inner:   h.inner.WithAttrs(attrs),
		newfunc: h.newfunc,
		opts:    h.opts,
		mu:      h.mu,
	}
}

func (h *OptHandler) WithGroup(name string) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	return &OptHandler{
		inner:   h.inner.WithGroup(name),
		newfunc: h.newfunc,
		opts:    h.opts,
		mu:      h.mu,
	}
}

func renew(h *OptHandler) {
	opts := h.opts
	h.inner = h.newfunc(&opts)
}

func (h *OptHandler) AddSource(addSource bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.opts.AddSource = addSource
	renew(h)
}

func (h *OptHandler) Level(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.opts.Level = level
	renew(h)
}

func (h *OptHandler) ReplaceAttr(replaceAttrr func(groups []string, a slog.Attr) slog.Attr) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.opts.ReplaceAttr = replaceAttrr
	renew(h)
}
