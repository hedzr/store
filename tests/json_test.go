package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/file"
)

func TestStore_JSON_Load(t *testing.T) {
	s := newBasicStore()
	if _, err := s.Load(
		context.TODO(),
		store.WithStorePrefix("app.json"),
		store.WithCodec(json.New()),
		store.WithProvider(file.New("../testdata/4.json")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.Equal(t, `r.Header.Get("From")`, s.MustGet("app.json.messages.0.placeholders.0.expr"))
	assert.Equal(t, `r.Header.Get("User-Agent")`, s.MustGet("app.json.messages.1.placeholders.0.expr"))
}
