package test_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/nestext"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := nestext.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.nested-text"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/7.txt")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `board member`, s.MustGet("app.nested-text.president.additional-roles.0"))
	assert.Equal(t, `1-210-555-8470`, s.MustGet("app.nested-text.Katheryn-McDaniel.phone.home"))
}

func TestNew2(t *testing.T) {
	s := store.New()
	parser := nestext.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.nested-text"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/8.txt")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `rsync borgbase`, s.MustGet("app.nested-text.repositories.home.children"))
}
