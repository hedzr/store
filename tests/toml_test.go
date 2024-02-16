package tests_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/providers/file"
)

func TestTOML(t *testing.T) {
	s := store.New()
	parser := toml.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.toml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/5.toml")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `127.0.0.1`, s.MustGet("app.toml.host"))
	assert.Equal(t, `TLS 1.3`, s.MustGet("app.toml.TLS.version"))
	assert.Equal(t, `AEAD-AES128-GCM-SHA256`, s.MustGet("app.toml.TLS.cipher"))
	assert.Equal(t, `go`, s.MustGet("app.toml.tags.0"))
}
