Package shandler provides some slog.Handlers.

[![](https://godoc.org/github.com/shu-go/shandler?status.svg)](https://godoc.org/github.com/shu-go/shandler)
[![Go Report Card](https://goreportcard.com/badge/github.com/shu-go/shandler)](https://goreportcard.com/report/github.com/shu-go/shandler)
![MIT License](https://img.shields.io/badge/License-MIT-blue)

# Examples

```go
import (
	"github.com/shu-go/shandler/leveled"
	"github.com/shu-go/shandler/multi"
	"github.com/shu-go/shandler/opt"

	fatihsan "github.com/fatih/color"
	"github.com/shu-go/shandler/color"
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

func Example_color() {
	scheme := color.DefaultDarkScheme()
	scheme.Level[slog.LevelDebug] = color.NewColor(fatihsan.FgHiWhite, fatihsan.BlinkRapid)
	scheme.Message = color.NewColor(fatihsan.Underline, fatihsan.Bold)

	h := color.NewHandler(os.Stdout, nil, scheme)
	slog.SetDefault(slog.New(h))

	slog.With(slog.String("s0", "value1")).WithGroup("grp1").With(slog.Int("i1", 1)).WithGroup("grp2").Info(
		"message",
		slog.String("str1", "value1"),
		slog.Int("int2", 2),
	)
}
```

----

Copyright 2024 Shuhei Kubota

<!--  vim: set et ft=markdown sts=4 sw=4 ts=4 tw=0 : -->
