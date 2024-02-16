package tests_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
)

func TestStore(t *testing.T) {
	t.Log("")
}

func TestStore_Get(t *testing.T) {
	conf := newBasicStore()

	// ret := conf.dump()
	// t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, conf.Search("apple"), `expecting conf.Search("apple") return true`)     // 返回 True
	// assertFalse(t, conf.Search("app"), `expecting conf.Search("app") return false`)       // 返回 False
	// assertTrue(t, conf.StartsWith("app"), `expecting conf.StartsWith("app") return true`) // 返回 True
	// conf.Insert("app")
	// assertTrue(t, conf.Search("app"), `expecting conf.Search("app") return true`) // 返回 True

	for i, c := range []struct {
		query string
		found bool
		data  any
	}{
		{"app.debug", true, false},
		{"app.dump", true, 3},
		{"app.d", false, nil},
		{"app.de", false, nil},
		{"app.deb", false, nil},
		{"app.verbose", true, true},
		{"app.logging", true, nil},
		{"app.logging.", true, nil},
		{"app.logging.file", true, "/tmp/1.log"},
		{"app.logging.rotate", true, 6},
		{"app.logging.words", true, []string{"a", "1", "false"}},
		{"app.server.start", true, 5},
		{"app", true, nil},
		{"app.", true, nil},
		{"a", false, nil},
	} {
		if data, found := conf.Get(c.query); found == c.found && reflect.DeepEqual(data, c.data) {
			continue
		} else { //nolint:revive
			t.Fatalf("%5d. querying %q and got (%v, %v), but expecting (%v, %v)",
				i, c.query, found, data, c.found, c.data)
		}
	}
}

func TestStore_Set(t *testing.T) {
	// trie := NewStoreT[any]()
	trie := newBasicStore()
	trie.Set("app.debug.map", map[string]any{"tags": []string{"delve", "verbose"}, "verbose": true})
	t.Logf("\nPath\n%v\n", trie.Dump())
}

func TestStore_Merge(t *testing.T) {
	// trie := NewStoreT[any]()
	trie := newBasicStore()
	if err := trie.Merge("app.debug.map", map[string]any{"tags": []string{"delve", "verbose"}, "verbose": true}); err != nil {
		t.Fatalf(`Merge failed: %v`, err)
	}
	t.Logf("\nPath\n%v\n", trie.Dump())

	assert.Equal(t, true, slices.Equal(trie.MustGet("app.debug.map.tags").([]string), []string{"delve", "verbose"}),
		`expecting trie.Get("app.debug.map.tags") return '[delve verbose]'`)
}

func TestStore_WithPrefix(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	assert.Equal(t, 6, trie.MustGet("app.logging.rotate"))
}

func TestStore_Clone(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	newStore := trie.Dup()
	t.Logf("\n[newStore] Path\n%v\n", newStore.Dump())

	assert.Equal(t, 6, newStore.MustGet("app.logging.rotate"))
}

func TestStore_GetR(t *testing.T) {
	trie := newBasicStore()
	trie.Set("app.dump.to", "stdout")
	t.Logf("\nPath\n%v\n", trie.Dump())

	ret := trie.MustR("app")
	t.Logf("\n[newStore] R (map)\n%v\n", ret)
	t.Logf("\n[newStore] R (map)\n")
	spew.Println(ret)

	t.Logf("Get: %s => %v", "app.debug", trie.MustFloat64("app.debug"))
	t.Logf("Get: %s => %v", "app.deb", trie.MustFloat64("app.deb")) // bad key test
	t.Logf("Get: %s => %v", "app.d", trie.MustFloat64("app.d"))     // bad key test
	t.Logf("Get: %s => %v", "app.", trie.MustFloat64("app."))       // bad key test

	t.Logf("Get: %s => %v", "app.dump.", trie.MustFloat64("app.dump."))
	t.Logf("Get: %s => %v", "app.dump", trie.MustFloat64("app.dump"))
	assert.Equal(t, 3.0, trie.MustFloat64("app.dump"))
	assert.Equal(t, 0.0, trie.MustFloat64("app.dump."))

	_, err := trie.GetFloat64("app.dump.")
	t.Logf("Get 'app.dump.', should be notfound: %v | notfound: %v", err, store.ErrIsNotFound(err))

	// assert.EqualTrue(t, newStore.MustGet("app.logging.rotate") == 6, `expecting trie.Get("app.logging.rotate") return 6`)
}

func newBasicStore() store.Store {
	conf := store.New()
	conf.Set("app.debug", false)
	// t.Logf("\nPath\n%v\n", conf.dump())
	conf.Set("app.verbose", true)
	// t.Logf("\nPath\n%v\n", conf.dump())
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	// conf.Set("app.logging.rotate", 6)
	// conf.Set("app.logging.words", []string{"a", "1", "false"})

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	return conf
}
