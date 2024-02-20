package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/maps"
)

func TestStore_maps_Load(t *testing.T) {
	s := newBasicStore()
	m := newMap()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.maps"),
		// store.WithCodec(json.New()),
		store.WithProvider(maps.New(m, "_")),

		store.WithStoreFlattenSlice(true), // expand map or slice in value
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.Equal(t, true, s.MustGet("app.maps.cool.station.8s"))
	assert.Equal(t, 5*time.Hour, s.MustGet("app.maps.cool.station.flush.interval"))
}
