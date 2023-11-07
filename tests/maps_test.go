package tests

import (
	"testing"
	"time"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/maps"
)

func TestStore_maps_Load(t *testing.T) {
	s := newBasicStore()
	m := newMap()
	if err := s.Load(
		store.WithStorePrefix("app.maps"),
		// store.WithCodec(json.New()),
		store.WithProvider(maps.New(m, "_")),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.EqualTrue(t, s.MustGet("app.maps.cool.station.8s") == true, `expecting store.Get("app.maps.cool.station.8s") return 'true'`)
	assert.EqualTrue(t, s.MustGet("app.maps.cool.station.flush.interval") == 5*time.Hour, `expecting store.Get("app.maps.cool.station.flush.interval") return '5h0m0s'`)
}

func newMap() map[string]any {
	return map[string]any{
		"cool.station": map[any]any{
			8 * time.Second: true,
			"flush": map[string]any{
				"always":   false,
				"interval": 5 * time.Hour,
			},
		},
		"desc": "a desc string",
	}
}
