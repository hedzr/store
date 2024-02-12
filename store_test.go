package store

import (
	"fmt"
	"reflect"
	"slices"
	"testing"

	"github.com/hedzr/store/ctx"
)

func TestCtxS_NamesCount(t *testing.T) {
	ctx := ctx.TODO()
	for ctx.Next() {
		t.Log(ctx.Key())
	}
}

func TestCtxS_Next(t *testing.T) {
	ctx := ctx.WithValue(ctx.TODO(), "k1", 1)
	ctx.WithValues("k2", 2, "k3", 3)

	for ctx.Next() {
		t.Log(ctx.Key())
	}
}

func TestStore_Get(t *testing.T) {
	trie := NewStoreT[any]()
	trie.Set("app.debug", 1)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Set("app.verbose", 2)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Set("app.dump", 3)
	trie.Set("app.logging.file", 4)
	trie.Set("app.server.start", 5)
	trie.Set("app.logging.rotate", 6)

	// ret := trie.dump()
	// t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	for i, c := range []struct {
		query string
		found bool
		data  any
	}{
		{"app.debug", true, 1},
		{"app.dump", true, 3},
		{"app.d", false, nil},
		{"app.de", false, nil},
		{"app.deb", false, nil},
		{"app.verbose", true, 2},
		{"app.logging", true, nil},
		{"app.logging.", true, nil},
		{"app.logging.file", true, 4},
		{"app.logging.rotate", true, 6},
		{"app.server.start", true, 5},
		{"app", true, nil},
		{"app.", true, nil},
		{"a", false, nil},
	} {
		if data, found := trie.Get(c.query); found == c.found && data == c.data {
			continue
		} else {
			t.Fatalf("%5d. querying %q and got (%v, %v), but expecting (%v, %v)", i, c.query, found, data, c.found, c.data)
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

	assertTrue(t, slices.Equal(trie.MustGet("app.debug.map.tags").([]string), []string{"delve", "verbose"}), `expecting trie.Get("app.debug.map.tags") return '[delve verbose]'`)
}

func TestStore_WithPrefix(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	assertTrue(t, trie.MustGet("app.logging.rotate") == 6, `expecting trie.Get("app.logging.rotate") return 6`)
}

func TestStore_Get2(t *testing.T) {
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
		} else {
			t.Fatalf("%5d. querying %q and got (%v, %v), but expecting (%v, %v)",
				i, c.query, found, data, c.found, c.data)
		}
	}
}

func assertTrue(t testing.TB, cond bool, msg ...any) { //nolint:govet //it's a printf/println dual interface
	if cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertTrue failed: %s", mesg)
}

func assertFalse(t testing.TB, cond bool, msg ...any) {
	if !cond {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertFalse failed: %s", mesg)
}

func newBasicStore() *storeS {
	conf := New()
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
