package tests_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hjson"
	"github.com/hedzr/store/providers/file"
)

func TestHjson(t *testing.T) {
	s := store.New()
	parser := hjson.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.hjson"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/6.hjson")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `r.Header.Get("From")`, s.MustGet("app.hjson.messages.0.placeholders.0.expr"))
	assert.Equal(t, `r.Header.Get("User-Agent")`, s.MustGet("app.hjson.messages.1.placeholders.0.expr"))
}
