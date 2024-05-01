package multi_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"

	"github.com/shu-go/gotwant"
	"github.com/shu-go/shandler/multi"
)

func TestMulti(t *testing.T) {
	buf1 := bytes.Buffer{}
	buf2 := bytes.Buffer{}
	buf3 := bytes.Buffer{}
	h := multi.NewHandler(
		slog.NewTextHandler(&buf1, nil),
		slog.NewTextHandler(&buf2, nil),
		slog.NewTextHandler(&buf3, nil),
	)
	slog.SetDefault(slog.New(h))

	slog.Info("hoge")

	gotwant.TestExpr(t, buf1.String(), strings.Contains(buf1.String(), "hoge"))
	gotwant.TestExpr(t, buf2.String(), strings.Contains(buf2.String(), "hoge"))
	gotwant.TestExpr(t, buf3.String(), strings.Contains(buf3.String(), "hoge"))
}
