package tests

import (
	"reflect"
	"slices"
	"testing"

	"github.com/hedzr/env/assert"
	"github.com/hedzr/env/assert/spew"
	"github.com/hedzr/store"
)

func TestStore_Get2(t *testing.T) {
	store := newBasicStore()

	// ret := store.dump()
	// t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, store.Search("apple"), `expecting store.Search("apple") return true`)     // 返回 True
	// assertFalse(t, store.Search("app"), `expecting store.Search("app") return false`)       // 返回 False
	// assertTrue(t, store.StartsWith("app"), `expecting store.StartsWith("app") return true`) // 返回 True
	// store.Insert("app")
	// assertTrue(t, store.Search("app"), `expecting store.Search("app") return true`) // 返回 True

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
		if data, found := store.Get(c.query); found == c.found && reflect.DeepEqual(data, c.data) {
			continue
		} else {
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

	assert.EqualTrue(t, slices.Equal(trie.MustGet("app.debug.map.tags").([]string), []string{"delve", "verbose"}), `expecting trie.Get("app.debug.map.tags") return '[delve verbose]'`)
}

func TestStore_WithPrefix(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	assert.EqualTrue(t, trie.MustGet("app.logging.rotate") == 6, `expecting newStore.Get("app.logging.rotate") return 6`)
}

func TestStore_Clone(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	newStore := trie.Dup()
	t.Logf("\n[newStore] Path\n%v\n", newStore.Dump())

	assert.EqualTrue(t, newStore.MustGet("app.logging.rotate") == 6, `expecting trie.Get("app.logging.rotate") return 6`)
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
	assert.EqualTrue(t, trie.MustFloat64("app.dump") == 3.0)
	assert.EqualTrue(t, trie.MustFloat64("app.dump.") == 0.0)

	_, err := trie.GetFloat64("app.dump.")
	t.Logf("Get 'app.dump.', should be notfound: %v | notfound: %v", err, store.ErrIsNotFound(err))

	// assert.EqualTrue(t, newStore.MustGet("app.logging.rotate") == 6, `expecting trie.Get("app.logging.rotate") return 6`)
}

func newBasicStore() store.Store {
	store := store.New()
	store.Set("app.debug", false)
	// t.Logf("\nPath\n%v\n", store.dump())
	store.Set("app.verbose", true)
	// t.Logf("\nPath\n%v\n", store.dump())
	store.Set("app.dump", 3)
	store.Set("app.logging.file", "/tmp/1.log")
	store.Set("app.server.start", 5)

	// store.Set("app.logging.rotate", 6)
	// store.Set("app.logging.words", []string{"a", "1", "false"})

	ss := store.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	return store
}
