package json_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := json.New()
	if err := s.Load(
		store.WithStorePrefix("app.json"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/4.json")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.json.messages.0.placeholders.0.expr") == `r.Header.Get("From")`, `expecting store.Get("app.json.messages.0.placeholders.0.expr") return 'r.Header.Get("From")'`)
	assert.EqualTrue(t, s.MustGet("app.json.messages.1.placeholders.0.expr") == `r.Header.Get("User-Agent")`, `expecting store.Get("app.json.messages.1.placeholders.0.expr") return 'r.Header.Get("User-Agent")'`)

}
