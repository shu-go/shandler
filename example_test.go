package shandler_test

import (
	"log/slog"
	"os"

	"github.com/shu-go/shandler/leveled"
	"github.com/shu-go/shandler/multi"
	"github.com/shu-go/shandler/opt"
)

func Example_multi() {
	h := multi.NewHandler(
		slog.NewTextHandler(os.Stdout, nil),
		slog.NewTextHandler(os.Stdout, nil),
	)
	slog.SetDefault(slog.New(h))

	slog.Info("one")
	slog.Info("two")
}

func Example_leveled() {
	h := leveled.NewHandler(
		slog.NewTextHandler(os.Stdout, nil),
		leveled.Debug(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})),
	)
	slog.SetDefault(slog.New(h))

	slog.Info("one")  // -> Default handler
	slog.Debug("two") // -> Debug handler
}

func Example_opt() {
	h := opt.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(h))

	slog.Info("one")  // output
	slog.Debug("two") // NO output

	h.Level(slog.LevelDebug)
	h.AddSource(true)

	slog.Info("three") // output with source
	slog.Debug("four") // output with source
}
