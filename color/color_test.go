package color_test

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"testing"
	"testing/slogtest"

	fatihsan "github.com/fatih/color"
	"github.com/shu-go/gotwant"
	"github.com/shu-go/shandler/color"
	stesting "github.com/shu-go/shandler/testing"
)

type logbackup struct {
	logger *slog.Logger
	flags  int
}

func backup() logbackup {
	return logbackup{
		logger: slog.Default(),
		flags:  log.Flags(),
	}
}

func (b logbackup) restore() {
	slog.SetDefault(b.logger)
	log.SetFlags(b.flags)
}

func TestColor(t *testing.T) {
	cb := &bytes.Buffer{}
	cl := slog.New(color.NewHandler(cb, nil, color.DefaultNilScheme()))
	//sb := &bytes.Buffer{}
	//sl := slog.New(slog.NewTextHandler(sb, nil))

	defer backup().restore()
	log.SetFlags(0)

	t.Run("Message", func(t *testing.T) {
		cb.Reset()
		cl.Info("message")
		gotwant.Test(t, cb.String(), "INFO message\n", gotwant.Format("%q"))

		cb.Reset()
		cl.Info("message 2")
		gotwant.Test(t, cb.String(), "INFO message 2\n", gotwant.Format("%q"))
	})

	t.Run("Attrs", func(t *testing.T) {
		cb.Reset()
		cl.Info("message", slog.String("str1", "value1"), slog.Int("int2", 2))
		gotwant.Test(t, cb.String(), "INFO message str1=value1 int2=2\n", gotwant.Format("%q"))

		cb.Reset()
		cl.Info("message", slog.Group("grp1", slog.String("str1", "value1"), slog.Group("grp2", slog.Int("int2", 2))))
		gotwant.Test(t, cb.String(), "INFO message grp1.str1=value1 grp1.grp2.int2=2\n", gotwant.Format("%q"))

		slog.Info("message", slog.Group("grp1", slog.String("str1", "value1"), slog.Int("int2", 2)))
	})

	t.Run("WithGroup", func(t *testing.T) {
		cb.Reset()
		l := cl.WithGroup("grp1")
		l.Info("message", slog.String("str1", "value1"), slog.Int("int2", 2))
		gotwant.Test(t, cb.String(), "INFO message grp1.str1=value1 grp1.int2=2\n", gotwant.Format("%q"))

		cb.Reset()
		l = l.WithGroup("grp2")
		l.Info("message", slog.String("str1", "value1"), slog.Int("int2", 2))
		gotwant.Test(t, cb.String(), "INFO message grp1.grp2.str1=value1 grp1.grp2.int2=2\n", gotwant.Format("%q"))

		slog.With().WithGroup("grp1").WithGroup("grp2").Info("message", slog.String("str1", "value1"), slog.Int("int2", 2))
	})

	t.Run("WithAttrsGroup", func(t *testing.T) {
		cb.Reset()
		cl.With(slog.String("s0", "value1")).WithGroup("grp1").With(slog.Int("i1", 1)).WithGroup("grp2").Info(
			"message",
			slog.String("str1", "value1"),
			slog.Int("int2", 2),
		)
		gotwant.Test(t, cb.String(), "INFO message s0=value1 grp1.i1=1 grp1.grp2.str1=value1 grp1.grp2.int2=2\n", gotwant.Format("%q"))

		slog.With(slog.String("s0", "value1")).WithGroup("grp1").With(slog.Int("i1", 1)).WithGroup("grp2").Info(
			"message",
			slog.String("str1", "value1"),
			slog.Int("int2", 2),
		)
	})

	t.Run("ReplaceAttr", func(t *testing.T) {
		cb.Reset()
		cl := slog.New(color.NewHandler(cb, &color.HandlerOptions{
			ReplaceAttr: func(group []string, a slog.Attr) slog.Attr {
				fmt.Fprintf(os.Stderr, "group: %+v\n", group)
				if strings.HasPrefix(a.Key, "str") {
					a.Value = slog.StringValue(strings.ToUpper(a.Value.String()))
				} else if strings.HasPrefix(a.Key, "int") {
					// groups: ["group1", "group2"]
					a.Value = slog.Int64Value(a.Value.Int64() * int64(len(group)+1))
				} else if strings.HasPrefix(a.Key, "bool") {
					a = slog.Attr{}
				}
				return a
			},
		}, color.DefaultNilScheme()))
		cl.With(slog.String("str0", "value1"), slog.Bool("bool1", true)).WithGroup("grp1").With(slog.Int("int1", 1)).WithGroup("grp2").Info(
			"message",
			slog.String("str1", "value1"),
			slog.Int("int2", 2),
			slog.Bool("bool2", false),
		)
		gotwant.Test(t, cb.String(), "INFO message str0=VALUE1 grp1.int1=2 grp1.grp2.str1=VALUE1 grp1.grp2.int2=6\n", gotwant.Format("%q"))

	})

	t.Run("slogtest", func(t *testing.T) {
		defer backup().restore()
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		cb.Reset()
		h := color.NewHandler(cb, &color.HandlerOptions{Compat: true}, color.DefaultNilScheme())

		err := slogtest.TestHandler(h, func() []map[string]any {
			return stesting.ParseTextLogs(t, cb.Bytes(), true)
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func TestColorShowcase(t *testing.T) {
	defer backup().restore()

	scheme := color.DefaultDarkScheme()
	scheme.Level[slog.LevelDebug] = color.NewColor(fatihsan.FgHiWhite, fatihsan.BlinkRapid)

	defer backup().restore()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	h := color.NewHandler(os.Stderr, &color.HandlerOptions{
		//AddSource: true,
		Level: slog.LevelDebug,
	}, scheme)
	l := slog.New(h)

	l = l.With("a0-1", "v0")
	l = l.With("a0-2", "v0")

	l = l.WithGroup("grp1")
	l.Info("info message", "a1", "v1", "a2", "value 2")
	l.Warn("warning message", "a1", "v1", "a2", "value 2")
	l.Error("error message", "a1", "v1", "a2", "value 2")
	l.Debug("debug message", "a1", "v1", "a2", "value 2")

	std := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	std = std.With("a0-1", "v0")
	std = std.With("a0-2", "v0")

	std = std.WithGroup("grp1")
	std.Info("info message", "a1", "v1", "a2", "value 2")
	std.Info("message", "a1", "v1", "a2", "value 2")

}

func BenchmarkColor(b *testing.B) {
	cb := &bytes.Buffer{}
	cl := slog.New(color.NewHandler(cb, nil, color.DefaultNilScheme()))

	b.Run("Simple", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cl.Info(
				"message",
				slog.String("str1", "value1"),
				slog.Int("int2", 2),
			)
		}
	})

	b.Run("WithGroup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cl.With(slog.String("s0", "value1")).WithGroup("grp1").With(slog.Int("i1", 1)).WithGroup("grp2").Info(
				"message",
				slog.String("str1", "value1"),
				slog.Int("int2", 2),
			)
		}
	})
}

func BenchmarkStd(b *testing.B) {
	cb := &bytes.Buffer{}
	cl := slog.New(slog.NewTextHandler(cb, nil))

	b.Run("Simple", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cl.Info(
				"message",
				slog.String("str1", "value1"),
				slog.Int("int2", 2),
			)
		}
	})

	b.Run("WithGroup", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cb.Reset()
			cl.With(slog.String("s0", "value1")).WithGroup("grp1").With(slog.Int("i1", 1)).WithGroup("grp2").Info(
				"message",
				slog.String("str1", "value1"),
				slog.Int("int2", 2),
			)
		}
	})
}
