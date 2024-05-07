package testing_test

import (
	"bytes"
	"log/slog"
	"testing"
	"testing/slogtest"

	stesting "github.com/shu-go/shandler/testing"
)

func TestTesting(t *testing.T) {
	t.Run("Text", func(t *testing.T) {
		buf := &bytes.Buffer{}
		h := slog.NewTextHandler(buf, nil)

		err := slogtest.TestHandler(h, func() []map[string]any {
			return stesting.ParseTextLogs(t, buf.Bytes(), true)
		})
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("JSON", func(t *testing.T) {
		buf := &bytes.Buffer{}
		h := slog.NewJSONHandler(buf, nil)

		err := slogtest.TestHandler(h, func() []map[string]any {
			return stesting.ParseJSONLogs(t, buf.Bytes(), true)
		})
		if err != nil {
			t.Error(err)
		}
	})
}
