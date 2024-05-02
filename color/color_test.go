package color_test

import (
	"bytes"
	"log"
	"log/slog"
	"os"
	"testing"

	fatihsan "github.com/fatih/color"
	"github.com/shu-go/gotwant"
	"github.com/shu-go/shandler/color"
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
}

func TestColorShowcase(t *testing.T) {
	defer backup().restore()

	scheme := color.DefaultDarkScheme()
	scheme.Level[slog.LevelDebug] = color.NewColor(fatihsan.FgHiWhite, fatihsan.BlinkRapid)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	h := color.NewHandler(os.Stderr, &slog.HandlerOptions{
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
