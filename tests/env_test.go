package tests

import (
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/env"
)

func TestStore_Env_Load(t *testing.T) {
	s := newBasicStore()
	if err := s.Load(
		store.WithStorePrefix("app.env"),
		store.WithProvider(env.New(
			env.WithPrefix(""),
			env.WithLowerCase(true),
			env.WithUnderlineToDot(true),
		)),
	); err != nil {
		t.Fatalf("LoadEnvTo failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.EqualTrue(t, s.MustGet("app.env.go111module") == `on`, `expecting store.Get("app.env.GO111MODULE") return 'on'`)

	// assertTrue(t, store.MustGet("app.env.HOME") == `/Users/hz`, `expecting store.Get("app.env.HOME") return '/Users/hz', but got '%`+`v'`, store.MustGet("app.env.HOME"))
}
