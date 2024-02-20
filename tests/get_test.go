package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/internal/times"
	"github.com/hedzr/store/providers/file"
)

func TestStoreResultTypedGetters(t *testing.T) {
	s := store.New()
	parser := hcl.New(hcl.WithFlattenSlices(true))
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.hcl"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.hcl")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, false, s.MustGet("app.hcl.debug").(bool))

	assert.Equal(t, `12345`, s.MustString("app.hcl.server.0.getters.0.mf"))
	assert.Equal(t, `-12.7831`, s.MustString("app.hcl.server.0.getters.0.ff"))
	assert.Equal(t, `true`, s.MustString("app.hcl.server.0.getters.0.bf"))
	assert.Equal(t, `1s52ms`, s.MustString("app.hcl.server.0.getters.0.tf"))
	assert.Equal(t, `2023-1-2`, s.MustGet("app.hcl.server.0.getters.0.time"))

	assert.Equal(t, int64(-12), s.MustInt64("app.hcl.server.0.getters.0.ff"))
}

func TestTypedGetters(t *testing.T) {
	s := store.New()
	parser := toml.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.toml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.toml")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, int64(5000), s.MustGet("app.toml.database.connection_max"))

	assert.Equal(t, times.MustSmartParseTime("1979-05-27 07:32:00 -0800"), s.MustGet("app.toml.owner.dob"))
	assert.Equal(t, time.Minute+37*time.Second+512*time.Millisecond, times.MustParseDuration(s.MustString("app.toml.owner.duration")))
}

func TestDecompoundMap(t *testing.T) {
	conf := newBasicStore()

	conf.Set("app.map", false)
	err := conf.Merge("app.map", map[string]any{
		"k1": 1,
		"k2": false,
		"m3": map[string]any{
			"bobo": "joe",
		},
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	assert.Equal(t, int(1), conf.MustGet("app.map.k1"))
	assert.Equal(t, false, conf.MustGet("app.map.k2"))
	assert.Equal(t, "joe", conf.MustGet("app.map.m3.bobo"))
}

func TestTagAndComment(t *testing.T) {
	conf := newBasicStore()

	found := conf.Has("app.logging.rotate")
	println(found)
	node, isBranch, isPartialMatched, found := conf.Locate("app.logging.rotate")
	t.Logf("%v    | %v, %v, found: %v", node.Data(), isBranch, isPartialMatched, found)

	conf.Set("debug", false)
	conf.SetComment("debug", "a flag to identify app debug mode", "remarks here")
	conf.SetTag("debug", map[string]any{
		"handler": func() {},
	})

	node, _, _, found = conf.Locate("debug")
	if found {
		t.Log(node.Tag(), node.Description(), node.Comment())
	}

	println(conf.Dump())
}
