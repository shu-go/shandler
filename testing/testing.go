// Package testing provides log parsers for slogtest.
package testing

import (
	"bytes"
	"encoding/json"
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
		for {
			var key, value []byte
			line, key, value = nextTextComponent(line)
			if len(key) == 0 || len(value) == 0 {
				break
			}
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

func nextTextComponent(b []byte) (remaining, key, value []byte) {
	// key
	eqidx := bytes.IndexByte(b, '=')
	if eqidx == -1 {
		return nil, nil, nil
	}

	// value
	start := eqidx + 1
	if len(b) <= start {
		return nil, b[:eqidx], nil
	}
	if b[start] != '"' {
		bb := b[start:]
		spidx := bytes.IndexByte(bb, ' ')
		if spidx == -1 {
			return nil, b[:eqidx], bb
		}
		return bb[spidx+1:], b[:eqidx], bb[:spidx]
	}
	start++
	for {
		bb := b[start:]
		qtidx := bytes.IndexByte(bb, '"')
		btidx := bytes.IndexByte(bb, '\\')
		if btidx == qtidx-1 {
			start = qtidx + 1
			continue
		}

		if qtidx == -1 {
			return nil, b[:eqidx], b[eqidx+1:]
		}
		return b[start+qtidx+1:], b[:eqidx], b[eqidx+1 : start+qtidx]
	}
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
