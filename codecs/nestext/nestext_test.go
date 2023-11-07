package nestext_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/nestext"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := nestext.New()
	if err := s.Load(
		store.WithStorePrefix("app.nested-text"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/7.txt")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.nested-text.president.additional-roles.0") == `board member`, `expecting store.Get("app.nested-text.president.additional-roles.0") return 'board member'`)
	assert.EqualTrue(t, s.MustGet("app.nested-text.Katheryn-McDaniel.phone.home") == `1-210-555-8470`, `expecting store.Get("app.nested-text.Katheryn-McDaniel.phone.home") return '1-210-555-8470'`)

}

func TestNew2(t *testing.T) {
	s := store.New()
	parser := nestext.New()
	if err := s.Load(
		store.WithStorePrefix("app.nested-text"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/8.txt")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.nested-text.repositories.home.children") == `rsync borgbase`, `expecting store.Get("app.nested-text.repositories.home.children") return 'rsync borgbase'`)
}
