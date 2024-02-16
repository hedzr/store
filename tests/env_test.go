package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/env"
)

func newBasicStore() store.Store {
	conf := store.New()
	conf.Set("app.debug", false)
	// t.Logf("\nPath\n%v\n", store.dump())
	conf.Set("app.verbose", true)
	// t.Logf("\nPath\n%v\n", store.dump())
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	// store.Set("app.logging.rotate", 6)
	// store.Set("app.logging.words", []string{"a", "1", "false"})

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	return conf
}

func TestStore_Env_Load(t *testing.T) {
	s := newBasicStore()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.env"),
		store.WithProvider(env.New(
			env.WithPrefix(""),
			env.WithLowerCase(true),
			env.WithUnderlineToDot(true),
		)),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("LoadEnvTo failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assert.Equal(t, `on`, s.MustGet("app.env.go111module"))

	// assertTrue(t, store.MustGet("app.env.HOME") == `/Users/hz`,
	//    `expecting store.Get("app.env.HOME") return '/Users/hz', but got '%`+`v'`,
	//    store.MustGet("app.env.HOME"))
}
