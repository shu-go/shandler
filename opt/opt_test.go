package opt_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/shu-go/shandler/opt"
)

func TestTarget(t *testing.T) {
	h := opt.NewTextHandler(os.Stderr, nil)
	slog.SetDefault(slog.New(h))

	slog.Debug("one")
	h.Level(slog.LevelDebug)
	slog.Debug("two")
	h.Level(slog.LevelInfo)
	slog.Debug("three")
}
