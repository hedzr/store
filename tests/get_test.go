package tests_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/env/times"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/providers/file"
)

func TestStoreResultTypedGetters(t *testing.T) {
	s := store.New()
	parser := hcl.New(hcl.WithFlattenSlices(true))
	if err := s.Load(
		store.WithStorePrefix("app.hcl"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.hcl")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualFalse(t, s.MustGet("app.hcl.debug").(bool),
		`expecting store.Get("app.hcl.debug") return 'false'`)

	assert.EqualTrue(t, s.MustString("app.hcl.server.0.getters.0.mf") == `12345`,
		`expecting store.GetString("app.hcl.server.0.getters.0.mf") return '-12.7831'`)
	assert.EqualTrue(t, s.MustString("app.hcl.server.0.getters.0.ff") == `-12.7831`,
		`expecting store.GetString("app.hcl.server.0.getters.0.ff") return '-12.7831'`)
	assert.EqualTrue(t, s.MustString("app.hcl.server.0.getters.0.bf") == `true`,
		`expecting store.GetString("app.hcl.server.0.getters.0.bf") return 'true'`)
	assert.EqualTrue(t, s.MustString("app.hcl.server.0.getters.0.tf") == `1s52ms`,
		`expecting store.GetString("app.hcl.server.0.getters.0.tf") return '1s52ms'`)
	assert.EqualTrue(t, s.MustGet("app.hcl.server.0.getters.0.time") == `2023-1-2`,
		`expecting store.GetString("app.hcl.server.0.getters.0.time") return '2023-1-2'`)

	assert.EqualTrue(t, s.MustInt64("app.hcl.server.0.getters.0.ff") == -12,
		fmt.Sprintf(`expecting store.GetString("app.hcl.server.0.getters.0.ff") return '%v'`,
			s.MustInt64("app.hcl.server.0.getters.0.ff")))
}

func TestTypedGetters(t *testing.T) {
	s := store.New()
	parser := toml.New()
	if err := s.Load(
		store.WithStorePrefix("app.toml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.toml")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.EqualTrue(t, s.MustGet("app.toml.database.connection_max") == int64(5000),
		fmt.Sprintf(`expecting store.Get("app.toml.database.connection_max") return 5000, but got %t.`,
			s.MustGet("app.toml.database.connection_max")))

	assert.EqualTrue(t, s.MustGet("app.toml.owner.dob") == times.MustSmartParseTime("1979-05-27 07:32:00 -0800"),
		fmt.Sprintf(`expecting store.Get("app.toml.owner.dob") return '1979-05-27 07:32:00 -0800', but got %t.`,
			s.MustGet("app.toml.owner.dob")))
	assert.EqualTrue(t, times.MustParseDuration(s.MustString("app.toml.owner.duration")) == time.Minute+37*time.Second+512*time.Millisecond,
		fmt.Sprintf(`expecting store.Get("app.toml.owner.dob") return '1m37s512ms', but got %t.`,
			s.MustGet("app.toml.owner.duration")))
}
