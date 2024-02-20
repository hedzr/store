package test_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/yaml"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := yaml.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.yaml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/2.yaml")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `-s`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.0"))
	assert.Equal(t, `-w`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.1"))

}

func TestNew2(t *testing.T) {
	s := store.New()
	parser := yaml.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.yaml"),
		store.WithCodec(parser),
		store.WithProvider(
			file.New("../../../testdata/2.yaml",
				file.WithPosition("app.bgo"),
			)),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `-s`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.0"))
	assert.Equal(t, `-w`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.1"))
}
