package tests_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/providers/file"
)

func TestTOML(t *testing.T) {
	s := store.New()
	parser := toml.New()
	if err := s.Load(
		store.WithStorePrefix("app.toml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/5.toml")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.toml.host") == `127.0.0.1`, `expecting store.Get("app.toml.host") return '127.0.0.1'`)
	assert.EqualTrue(t, s.MustGet("app.toml.TLS.version") == `TLS 1.3`, `expecting store.Get("app.toml.TLS.version") return 'TLS 1.3'`)
	assert.EqualTrue(t, s.MustGet("app.toml.TLS.cipher") == `AEAD-AES128-GCM-SHA256`,
		`expecting store.Get("app.toml.TLS.cipher") return 'AEAD-AES128-GCM-SHA256'`)
	assert.EqualTrue(t, s.MustGet("app.toml.tags.0") == `go`, `expecting store.Get("app.toml.tags.0") return 'go'`)
}
