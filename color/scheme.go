package color

import (
	"log/slog"

	"github.com/fatih/color"
)

//////////////////////////////////////////////////

type Scheme struct {
	Base Colorizer

	Level     map[slog.Level]Colorizer
	Time      Colorizer
	Source    Colorizer
	Message   Colorizer
	AttrKey   Colorizer
	AttrValue Colorizer
}

func (s Scheme) LevelPrinter(level slog.Level) Colorizer {
	if lp, found := s.Level[level]; found {
		if lp != nil {
			return lp
		}
	}
	return s.BasePrinter()
}

func (s Scheme) TimePrinter() Colorizer {
	if s.Time != nil {
		return s.Time
	}
	return s.BasePrinter()
}

func (s Scheme) SourcePrinter() Colorizer {
	if s.Source != nil {
		return s.Source
	}
	return s.BasePrinter()
}

func (s Scheme) MessagePrinter() Colorizer {
	if s.Message != nil {
		return s.Message
	}
	return s.BasePrinter()
}

func (s Scheme) AttrKeyPrinter() Colorizer {
	if s.AttrKey != nil {
		return s.AttrKey
	}
	return s.BasePrinter()
}

func (s Scheme) AttrValuePrinter() Colorizer {
	if s.AttrValue != nil {
		return s.AttrValue
	}
	return s.BasePrinter()
}

func (s Scheme) BasePrinter() Colorizer {
	if s.Base != nil {
		return s.Base
	}
	return defaultColorizer
}

func DefaultNilScheme() *Scheme {
	return &Scheme{}
}

func DefaultLightScheme() *Scheme {
	return &Scheme{
		Base: NewColor(color.FgHiBlack, color.Faint),

		Message:   NewColor(color.FgBlack),
		AttrKey:   NewColor(color.FgBlue),
		AttrValue: NewColor(color.FgBlack),
		Level: map[slog.Level]Colorizer{
			slog.LevelInfo:  NewColor(color.FgHiBlack, color.Faint),
			slog.LevelWarn:  NewColor(color.FgYellow, color.Bold),
			slog.LevelError: NewColor(color.FgRed, color.Bold),
			slog.LevelDebug: NewColor(color.FgBlack),
		},
	}
}

func DefaultDarkScheme() *Scheme {
	return &Scheme{
		Base: NewColor(color.FgWhite, color.Faint),

		Message:   NewColor(color.FgHiWhite),
		AttrKey:   NewColor(color.FgHiCyan),
		AttrValue: NewColor(color.FgHiWhite),
		Level: map[slog.Level]Colorizer{
			slog.LevelInfo:  NewColor(color.FgWhite, color.Faint),
			slog.LevelWarn:  NewColor(color.FgYellow, color.Bold),
			slog.LevelError: NewColor(color.FgHiRed, color.Bold),
			slog.LevelDebug: NewColor(color.FgWhite, color.Faint),
		},
	}
}
