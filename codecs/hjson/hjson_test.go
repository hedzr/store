package hjson_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hjson"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := hjson.New()
	if err := s.Load(
		store.WithStorePrefix("app.hjson"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/6.hjson")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.hjson.messages.0.placeholders.0.expr") == `r.Header.Get("From")`, `expecting store.Get("app.hjson.messages.0.placeholders.0.expr") return 'r.Header.Get("From")'`)
	assert.EqualTrue(t, s.MustGet("app.hjson.messages.1.placeholders.0.expr") == `r.Header.Get("User-Agent")`, `expecting store.Get("app.hjson.messages.1.placeholders.0.expr") return 'r.Header.Get("User-Agent")'`)

}
