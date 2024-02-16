package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/yaml"
	"github.com/hedzr/store/providers/file"
)

func TestStore_YAML_Load(t *testing.T) {
	s := newBasicStore()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.yaml"),
		store.WithCodec(yaml.New()),
		store.WithProvider(file.New("../testdata/2.yaml")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)
	assert.Equal(t, `-s`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.0"))
	assert.Equal(t, `-w`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.1"))
}

func TestStore_YAML_Load_Normal(t *testing.T) {
	s := newBasicStore()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.yaml"),
		store.WithCodec(yaml.New()),
		store.WithProvider(file.New("../testdata/2.yaml")),

		// store.WithStoreFlattenSlice(false),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)
	assert.Equal(t, []interface{}{`-s`, "-w"}, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags"))
}
