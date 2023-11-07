package yaml_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/yaml"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := yaml.New()
	if err := s.Load(
		store.WithStorePrefix("app.yaml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/2.yaml")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.0") == `-s`, `expecting store.Get("001-bgo.001-bgo.ldflags.0") return '-s'`)
	assert.EqualTrue(t, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.1") == `-w`, `expecting store.Get("001-bgo.001-bgo.ldflags.1") return '-w'`)

}
