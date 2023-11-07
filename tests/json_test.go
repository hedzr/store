package tests

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/file"
)

func TestStore_JSON_Load(t *testing.T) {
	s := newBasicStore()
	if err := s.Load(
		store.WithStorePrefix("app.json"),
		store.WithCodec(json.New()),
		store.WithProvider(file.New("../testdata/4.json")),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.EqualTrue(t, s.MustGet("app.json.messages.0.placeholders.0.expr") == `r.Header.Get("From")`, `expecting store.Get("app.json.messages.0.placeholders.0.expr") return 'r.Header.Get("From")'`)
	assert.EqualTrue(t, s.MustGet("app.json.messages.1.placeholders.0.expr") == `r.Header.Get("User-Agent")`, `expecting store.Get("app.json.messages.1.placeholders.0.expr") return 'r.Header.Get("User-Agent")'`)
}
