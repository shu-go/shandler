package leveled_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/shu-go/gotwant"
	"github.com/shu-go/shandler/leveled"
)

func TestDefault(t *testing.T) {
	defaultBuf := bytes.Buffer{}

	h := leveled.NewHandler(
		slog.NewTextHandler(&defaultBuf, nil),
	)
	slog.SetDefault(slog.New(h))

	slog.Debug("one")
	slog.Error("two")
	slog.Info("three")
	slog.Warn("four")

	//gotwant.TestExpr(t, defaultBuf.String(), strings.Contains(defaultBuf.String(), "one"))
	gotwant.TestExpr(t, defaultBuf.String(), strings.Contains(defaultBuf.String(), "two"))
	gotwant.TestExpr(t, defaultBuf.String(), strings.Contains(defaultBuf.String(), "three"))
	gotwant.TestExpr(t, defaultBuf.String(), strings.Contains(defaultBuf.String(), "four"))
}

func TestSimple(t *testing.T) {
	defaultBuf := bytes.Buffer{}
	debugBuf := bytes.Buffer{}
	errorBuf := bytes.Buffer{}
	infoBuf := bytes.Buffer{}
	warnBuf := bytes.Buffer{}

	h := leveled.NewHandler(
		slog.NewTextHandler(&defaultBuf, nil),
		leveled.Debug(slog.NewTextHandler(&debugBuf, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
		leveled.Error(slog.NewTextHandler(&errorBuf, nil)),
		leveled.Info(slog.NewTextHandler(&infoBuf, nil)),
		leveled.Warn(slog.NewTextHandler(&warnBuf, nil)),
	)
	slog.SetDefault(slog.New(h))

	slog.Debug("one")
	slog.Error("two")
	slog.Info("three")
	slog.Warn("four")

	gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "one"))
	gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "two"))
	gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "three"))
	gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "four"))
}

func TestWith(t *testing.T) {
	t.Run("WithAttrs", func(t *testing.T) {
		defaultBuf := bytes.Buffer{}
		debugBuf := bytes.Buffer{}
		errorBuf := bytes.Buffer{}
		infoBuf := bytes.Buffer{}
		warnBuf := bytes.Buffer{}

		h := leveled.NewHandler(
			slog.NewTextHandler(&defaultBuf, nil),
			leveled.Debug(slog.NewTextHandler(&debugBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})),
			leveled.Error(slog.NewTextHandler(&errorBuf, nil)),
			leveled.Info(slog.NewTextHandler(&infoBuf, nil)),
			leveled.Warn(slog.NewTextHandler(&warnBuf, nil)),
		)
		slog.SetDefault(slog.New(h))

		wlog := slog.With("attr1", "value1")

		wlog.Debug("one")
		wlog.Error("two")
		wlog.Info("three")
		wlog.Warn("four")

		println(debugBuf.String())
		gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "one"))
		gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "attr1=value1"))

		gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "two"))
		gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "attr1=value1"))

		gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "three"))
		gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "attr1=value1"))

		gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "four"))
		gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "attr1=value1"))
	})

	t.Run("WithGroup", func(t *testing.T) {
		defaultBuf := bytes.Buffer{}
		debugBuf := bytes.Buffer{}
		errorBuf := bytes.Buffer{}
		infoBuf := bytes.Buffer{}
		warnBuf := bytes.Buffer{}

		h := leveled.NewHandler(
			slog.NewTextHandler(&defaultBuf, nil),
			leveled.Debug(slog.NewTextHandler(&debugBuf, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			})),
			leveled.Error(slog.NewTextHandler(&errorBuf, nil)),
			leveled.Info(slog.NewTextHandler(&infoBuf, nil)),
			leveled.Warn(slog.NewTextHandler(&warnBuf, nil)),
		)
		slog.SetDefault(slog.New(h))

		wlog := slog.With("attr1", "value1").WithGroup("group1").With("attr2", "value2")

		wlog.Debug("one")
		wlog.Error("two")
		wlog.Info("three")
		wlog.Warn("four")

		println(debugBuf.String())
		gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "one"))
		gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "attr1=value1"))
		gotwant.TestExpr(t, debugBuf.String(), strings.Contains(debugBuf.String(), "group1"))

		gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "two"))
		gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "attr1=value1"))
		gotwant.TestExpr(t, errorBuf.String(), strings.Contains(errorBuf.String(), "group1"))

		gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "three"))
		gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "attr1=value1"))
		gotwant.TestExpr(t, infoBuf.String(), strings.Contains(infoBuf.String(), "group1"))

		gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "four"))
		gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "attr1=value1"))
		gotwant.TestExpr(t, warnBuf.String(), strings.Contains(warnBuf.String(), "group1"))
	})
}
