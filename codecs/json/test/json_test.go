package test_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := json.New()
	if err := s.Load(context.TODO(),
		store.WithStorePrefix("app.json"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/4.json")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `r.Header.Get("From")`, s.MustGet("app.json.messages.0.placeholders.0.expr").(string))
	assert.Equal(t, `r.Header.Get("User-Agent")`, s.MustGet("app.json.messages.1.placeholders.0.expr").(string))
}
