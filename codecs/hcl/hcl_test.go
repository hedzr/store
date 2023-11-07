package hcl_test

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := hcl.New(hcl.WithFlattenSlices(true))
	if err := s.Load(
		store.WithStorePrefix("app.hcl"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../testdata/9.hcl")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.hcl.service.1.http.0.web_proxy.0.process.0.main.0.command.0") == `/usr/local/bin/awesome-app`, `expecting store.Get("app.hcl.service.1.http.0.web_proxy.0.process.0.main.0.command.0") return '/usr/local/bin/awesome-app'`)

}
