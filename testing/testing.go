// Package testing provides log parsers for slogtest.
package testing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	gotesting "testing"

	"golang.org/x/exp/maps"
)

func ParseJSONLogs(t *gotesting.T, in []byte, show bool) []map[string]any {
	t.Helper()

	ms := make([]map[string]any, 0)
	for i, line := range bytes.Split(in, []byte{'\n'}) {
		if len(line) == 0 {
			if show {
				t.Logf("%02d: (empty)\n", i+1)
			}
			continue
		}

		if show {
			t.Logf("%02d: %s\n", i+1, string(line))
		}

		m := make(map[string]any)
		err := json.Unmarshal(line, &m)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
		m = group(m)
		if show {
			t.Logf(" => %+v\n", m)
		}
		ms = append(ms, m)
	}
	return ms
}

func ParseTextLogs(t *gotesting.T, in []byte, show bool) []map[string]any {
	t.Helper()

	ms := make([]map[string]any, 0)
	for i, line := range bytes.Split(in, []byte{'\n'}) {
		if len(line) == 0 {
			if show {
				t.Logf("%02d: (empty)\n", i+1)
			}
			continue
		}

		if show {
			t.Logf("%02d: %s\n", i+1, string(line))
		}

		m := make(map[string]any)
		for _, c := range bytes.Split(line, []byte{' '}) {
			eqidx := bytes.Index(c, []byte{'='})
			if eqidx == -1 {
				panic(fmt.Sprintf("%02d: line:%q, component:%q, eqidx:%d", i+1, line, c, eqidx))
			}
			key := c[:eqidx]
			value := c[eqidx+1:]
			m[string(key)] = string(value)
		}
		m = group(m)
		if show {
			t.Logf(" => %+v\n", m)
		}
		ms = append(ms, m)
	}
	return ms
}

func group(m map[string]any) map[string]any {
	for {
		groupexists := false

		keys := maps.Keys(m)

		for _, k := range keys {
			lastidx := strings.LastIndex(k, ".")
			if lastidx == -1 {
				continue
			}

			groupexists = true

			v := m[k]
			origK := k
			g := k[:lastidx]
			k = k[lastidx+1:]

			if gm, found := m[g]; found {
				if gm, ok := gm.(map[string]any); ok {
					gm[k] = v
				}
			} else {
				gm := make(map[string]any)
				gm[k] = v
				m[g] = gm
			}
			delete(m, origK)
		}

		if !groupexists {
			break
		}
	}

	return m
}
