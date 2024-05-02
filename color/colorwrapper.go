package color

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
)

// Almost all colorizing codes are from "github.com/fatih/color"

type Colorizer interface {
	AppendFormat(b []byte) []byte
	AppendUnformat(b []byte) []byte
}

type fmtAppender struct{}

func (f fmtAppender) AppendFormat(b []byte) []byte {
	return b
}

func (f fmtAppender) AppendUnformat(b []byte) []byte {
	return b
}

var defaultColorizer = &fmtAppender{}

//////////////////////////////////////////////////

// Almost all codes are from "github.com/fatih/color"
type Color struct {
	*color.Color
	Params  []color.Attribute
	NoColor *bool
}

func NewColor(value ...color.Attribute) *Color {
	raw := color.New(value...)
	c := &Color{
		Color:  raw,
		Params: make([]color.Attribute, 0, len(value)),
	}
	if noColorIsSet() {
		c.NoColor = boolPtr(true)
	}
	c.Add(value...)
	return c
}

func (c *Color) Add(value ...color.Attribute) *Color {
	c.Color.Add(value...)
	c.Params = append(c.Params, value...)
	return c
}

func (c *Color) AppendFormat(b []byte) []byte {
	if c.isNoColorSet() {
		return b
	}

	b = append(b, '\x1b')
	b = append(b, '[')
	b = c.appendSequence(b)
	b = append(b, 'm')
	return b
}

func (c *Color) AppendUnformat(b []byte) []byte {
	if c.isNoColorSet() {
		return b
	}

	b = append(b, '\x1b')
	b = append(b, '[')

	for i, v := range c.Params {
		if i > 0 {
			b = append(b, ';')
		}

		b = fmt.Appendf(b, "%d", color.Reset)
		ra, ok := mapResetAttributes[v]
		if ok {
			b = fmt.Appendf(b, "%d", ra)
		}
	}

	b = append(b, 'm')
	return b
	//return fmt.Appendf(b, "%s[%sm", escape, strings.Join(format, ";"))
}

var mapResetAttributes = map[color.Attribute]color.Attribute{
	color.Bold:         color.ResetBold,
	color.Faint:        color.ResetBold,
	color.Italic:       color.ResetItalic,
	color.Underline:    color.ResetUnderline,
	color.BlinkSlow:    color.ResetBlinking,
	color.BlinkRapid:   color.ResetBlinking,
	color.ReverseVideo: color.ResetReversed,
	color.Concealed:    color.ResetConcealed,
	color.CrossedOut:   color.ResetCrossedOut,
}

func (c *Color) isNoColorSet() bool {
	// check first if we have user set action
	if c.NoColor != nil {
		return *c.NoColor
	}

	// if not return the global option, which is disabled by default
	return NoColor
}

var (
	NoColor = noColorIsSet() || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
)

func (c *Color) appendSequence(b []byte) []byte {
	for i, v := range c.Params {
		if i > 0 {
			b = append(b, ';')
		}
		b = fmt.Appendf(b, "%d", v)
	}
	return b
}

func noColorIsSet() bool {
	return os.Getenv("NO_COLOR") != ""
}

func boolPtr(v bool) *bool {
	return &v
}
