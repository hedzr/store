package test_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := hcl.New(hcl.WithFlattenSlices(true))
	if err := s.Load(context.TODO(),
		store.WithStorePrefix("app.hcl"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/9.hcl")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `/usr/local/bin/awesome-app`,
		s.MustGet("app.hcl.service.1.http.0.web_proxy.0.process.0.main.0.command.0").(string),
	)
	//	 `expecting store.Get("app.hcl.service.1.http.0.web_proxy.0.process.0.main.0.command.0") return '/usr/local/bin/awesome-app'`)
}
