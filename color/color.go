package color

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
)

type ColorHandler struct {
	scheme Scheme

	opts slog.HandlerOptions

	attrs  []byte
	groups []string

	mu *sync.Mutex
	w  io.Writer

	cacheFlagDate   bool
	cacheFlagTime   bool
	cacheFlagMicro  bool
	cacheTimeFormat string
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions, scheme *Scheme) *ColorHandler {
	h := &ColorHandler{
		mu: &sync.Mutex{},
		w:  w,
	}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}

	if scheme != nil {
		h.scheme = *scheme
	} else {
		scheme = DefaultDarkScheme()
		h.scheme = *scheme
	}

	return h
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := h.clone()

	prefix := strings.Join(h.groups, ".")

	pk := h.scheme.AttrKeyPrinter()
	pv := h.scheme.AttrValuePrinter()
	pn := h.scheme.BasePrinter()

	h2attrs := h2.attrs
	for i := 0; i < len(attrs); i++ {
		if attrs[i].Equal(slog.Attr{}) {
			continue
		}

		h2attrs = appendAttr(h2attrs, prefix, attrs[i], h.opts.ReplaceAttr, h.groups, pk, pv, pn)
	}

	h2.attrs = h2attrs

	return h2
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	pbuf := pool.Get().(*[]byte)
	buf := *pbuf
	buf = buf[:0]

	flags := log.Flags()
	flagdate := flags&log.Ldate != 0
	flagtime := flags&log.Ltime != 0
	flagmicro := flags&log.Lmicroseconds != 0
	flagshortfile := flags&log.Lshortfile != 0
	flaglongfile := flags&log.Llongfile != 0

	level := r.Level.Level()

	if !r.Time.IsZero() {
		var format string

		if flagdate == h.cacheFlagDate && flagtime == h.cacheFlagTime && flagmicro == h.cacheFlagMicro {
			format = h.cacheTimeFormat
		} else {
			if flagdate {
				format = "2006/01/02"
			}
			if flagtime || flagmicro {
				if format != "" {
					format += " "
				}

				format += "15:04:05"
				if flagmicro {
					format += ".000000"
				}
			}
			h.cacheFlagDate = flagdate
			h.cacheFlagTime = flagtime
			h.cacheFlagMicro = flagmicro
			h.cacheTimeFormat = format
		}

		if format != "" {
			tm := h.scheme.TimePrinter()
			buf = tm.AppendFormat(buf)
			if flags&log.LUTC != 0 {
				buf = r.Time.UTC().AppendFormat(buf, format)
			} else {
				buf = r.Time.AppendFormat(buf, format)
			}
			buf = tm.AppendUnformat(buf)
			buf = append(buf, ' ')
		}
	}

	if (h.opts.AddSource || flaglongfile || flagshortfile) && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()

		src := h.scheme.SourcePrinter()
		buf = src.AppendFormat(buf)
		if h.opts.AddSource || flaglongfile {
			buf = append(buf, f.File...)
			buf = append(buf, ':')
			buf = fmt.Appendf(buf, "%d", f.Line)
			buf = append(buf, ':')
			buf = append(buf, ' ')
		} else {
			buf = append(buf, filepath.Base(f.File)...)
			buf = append(buf, ':')
			buf = fmt.Appendf(buf, "%d", f.Line)
			buf = append(buf, ':')
			buf = append(buf, ' ')
		}
		buf = src.AppendUnformat(buf)
	}

	lvl := h.scheme.LevelPrinter(level)
	buf = lvl.AppendFormat(buf)
	buf = append(buf, level.String()...)
	buf = lvl.AppendUnformat(buf)

	buf = append(buf, ' ')

	msg := h.scheme.MessagePrinter()
	buf = msg.AppendFormat(buf)
	buf = append(buf, r.Message...)
	buf = msg.AppendUnformat(buf)

	buf = append(buf, h.attrs...)

	pk := h.scheme.AttrKeyPrinter()
	pv := h.scheme.AttrValuePrinter()
	pn := h.scheme.BasePrinter()

	r.Attrs(func(a slog.Attr) bool {
		if a.Equal(slog.Attr{}) {
			return true
		}

		prefix := strings.Join(h.groups, ".")
		buf = appendAttr(buf, prefix, a, h.opts.ReplaceAttr, h.groups, pk, pv, pn)

		return true
	})
	buf = append(buf, '\n')

	h.mu.Lock()
	_, err := h.w.Write(buf)
	h.mu.Unlock()

	pbuf = &buf
	pool.Put(pbuf)

	return err
}

func appendQuote(b []byte, s string) []byte {
	if strings.ContainsAny(s, " \t\"") {
		b = append(b, '"')
		b = append(b, s...)
		b = append(b, '"')
		return b
	}
	return append(b, s...)
}

func (h ColorHandler) clone() *ColorHandler {
	h2 := ColorHandler{
		opts:   h.opts,
		attrs:  slices.Clip(h.attrs),
		groups: slices.Clip(h.groups),
		mu:     h.mu,
		w:      h.w,
		scheme: h.scheme,
	}
	return &h2
}

func appendAttr(buf []byte, prefix string, a slog.Attr, rep func(groups []string, a slog.Attr) slog.Attr, groups []string, pk, pv, pb Colorizer) []byte {
	a.Value.Resolve()

	if a.Value.Kind() == slog.KindGroup {
		if prefix == "" {
			prefix = a.Key
		} else {
			prefix += "." + a.Key
		}
		for _, child := range a.Value.Group() {
			buf = appendAttr(buf, prefix, child, rep, groups, pk, pv, pb)
		}
	} else {
		if rep != nil {
			a = rep(groups, a)
			if a.Equal(slog.Attr{}) {
				return buf
			}
			a.Value = a.Value.Resolve()
		}

		buf = append(buf, ' ')

		buf = pk.AppendFormat(buf)
		if prefix != "" {
			buf = append(buf, prefix...)
			buf = append(buf, '.')
		}
		buf = append(buf, a.Key...)
		buf = pk.AppendUnformat(buf)

		buf = pb.AppendFormat(buf)
		buf = append(buf, '=')
		buf = pb.AppendUnformat(buf)

		buf = pv.AppendFormat(buf)
		buf = appendQuote(buf, a.Value.String())
		buf = pv.AppendUnformat(buf)
	}

	return buf
}

var pool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 256)
		return &b
	},
}
