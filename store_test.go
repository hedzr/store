package store

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/hedzr/store/radix"
)

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
		{"app.d", false, nil},   // partial-matched
		{"app.de", false, nil},  // non-exist key
		{"app.deb", false, nil}, // no-exist key
		{"app.verbose", true, 2},
		{"app.logging", true, nil},
		{"app.logging.", true, nil},
		{"app.logging.file", true, 4},
		{"app.logging.rotate", true, 6},
		{"app.server.start", true, 5},
		{"app", true, nil},  // a branch key
		{"app.", true, nil}, // a branch key with ending '.'
		{"a", false, nil},   // a missed key
	} {
		if data, found := trie.Get(c.query); found == c.found && data == c.data {
			continue
		} else {
			t.Fatalf("%5d. querying %q and got (%v, %v), but expecting (%v, %v)", i, c.query, found, data, c.found, c.data)
		}
	}
}

func TestStore_Get2(t *testing.T) { // test data with `reflect.DeepEqual`, tests for slice or map
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

func TestSimpleStore(t *testing.T) {
	conf := newBasicStore()
	conf.Set("app.logging.words", []any{"a", 1, false})
	conf.Set("app.server.sites", -1)
	t.Logf("\nPath\n%v\n", conf.Dump())
}

func TestStore_GetM(t *testing.T) {
	conf := newBasicStore()

	m, err := conf.GetM("")
	if err != nil {
		t.Fatalf("wrong in calling GetM(\"\"): %v", err)
	}
	t.Logf("whole tree: %v", m)

	// filter by functor

	m, err = conf.GetM("",
		WithKeepPrefix[any](false),
		WithoutFlattenKeys[any](false),
		WithFilter[any](func(node radix.Node[any]) bool {
			return strings.HasPrefix(node.Key(), "app.logging.")
		}))
	if err != nil {
		t.Fatalf("wrong in calling GetM(\"\"): %v", err)
	}
	t.Logf("app.logging sub-tree: %v", m)
}

func TestStore_GetSectionFrom(t *testing.T) {
	conf := newBasicStore()
	conf.Set("app.logging.words", []any{"a", 1, false})
	conf.Set("app.server.sites", -1)
	t.Logf("\nPath\n%v\n", conf.Dump())

	type loggingS struct {
		File   uint
		Rotate uint64
		Words  []any
	}

	type serverS struct {
		Start int
		Sites int
	}

	type appS struct {
		Debug   int
		Dump    int
		Verbose int64
		Logging loggingS
		Server  serverS
	}

	type cfgS struct {
		App appS
	}

	var ss cfgS
	err := conf.GetSectionFrom("", &ss)
	t.Logf("cfgS: %v | err: %v", ss, err)

	assertEqual(t, []any{"a", 1, false}, ss.App.Logging.Words)
	assertEqual(t, -1, ss.App.Server.Sites)

	if !reflect.DeepEqual(ss.App.Logging.Words, []any{"a", 1, false}) {
		t.Fail()
	}

	err = conf.GetSectionFrom("nonexist", nil)
	t.Log("nothing happened")

	err = conf.GetSectionFrom("nonexist", &ss)
	t.Logf("cfgS: %v | err: %v", ss, err)
	if err != nil {
		t.Fail()
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

	assertEqual(t, 6, trie.MustGet("app.logging.rotate"))
	conf := trie.WithPrefix("app")
	assertEqual(t, 6, conf.MustGet("logging.rotate"))
	conf = conf.WithPrefix("logging")
	assertEqual(t, 6, conf.MustGet("rotate"))
	conf = trie.WithPrefixReplaced("app.logging")
	assertEqual(t, 6, conf.MustGet("rotate"))
}

func TestStore_Walk(t *testing.T) {
	var conf Store = newBasicStore()
	conf.Walk("", func(path, fragment string, node radix.Node[any]) {
		t.Logf("%v / %v => %v", path, fragment, node)
	})
}

func TestStore_Walk1(t *testing.T) {
	var conf Store = newBasicStore()
	conf.Walk("app", func(path, fragment string, node radix.Node[any]) {
		t.Logf("%v / %v => %v", path, fragment, node)
	})
}

func TestStore_Walk2(t *testing.T) {
	var conf Store = newBasicStore(
		WithDelimiter('.'),
		WithPrefix(""),
		WithFlattenSlice(true),
		WithWatchEnable(false),
	)
	defer conf.Close()

	_ = conf.MustGet("x")
	_ = conf.Remove("x")
	_ = conf.Has("x")
	_, _, _, _ = conf.Locate("x")
	_ = conf.Clone()
	_ = conf.Dup()

	conf.SetPrefix(conf.Prefix())

	_ = conf.(*storeS).loadMap(map[string]any{
		"m1": map[string]any{
			"s1": "s1",
			"s2": 2,
		},
		"m2": []map[string]any{
			{
				"s1": "s1",
				"s2": 2,
			},
			{
				"s1": "s1",
				"s2": 2,
			},
		},
		"m3": []any{1, 2, false},
		"m4": 3.1313,
	}, "app", true, nil)

	conf.Walk("app.", func(path, fragment string, node radix.Node[any]) {
		t.Logf("%v / %v => %v", path, fragment, node)
	})
}

func TestStore_Monitoring(t *testing.T) {
	conf := newBasicStore(
		WithOnNewHandlers(func(path string, value any, mergingMapOrLoading bool) {
			t.Logf("[new] %q: %v | %v", path, value, mergingMapOrLoading)
		}),
		WithOnDeleteHandlers(func(path string, value any, mergingMapOrLoading bool) {
			t.Logf("[del] %q: %v | %v", path, value, mergingMapOrLoading)
		}),
		WithOnChangeHandlers(func(path string, value, oldValue any, mergingMapOrLoading bool) {
			t.Logf("[mod] %q: %v => %v | %v", path, oldValue, value, mergingMapOrLoading)
		}),
	)

	t.Logf("tests begins.")
	conf.Set("app.server.port", 7300)
	assertEqual(t, 7300, conf.MustInt("app.server.port", -1))
	assertTrue(t, 7300 == conf.MustInt("app.server.port", -1))
	conf.Set("app.server.port", 7301)
	assertEqual(t, 7301, conf.MustInt("app.server.port", -1))
	assertFalse(t, 7300 == conf.MustInt("app.server.port", -1))
	conf.Set("app.server.tls.cert", "/tmp/cert.pem")
	conf.Set("app.server.tls.priv", "/tmp/private-key.pem")
	assertEqual(t, "/tmp/private-key.pem", conf.MustString("app.server.tls.priv"))
	node, np, removed := conf.RemoveEx("app.server.tls.cert")
	t.Logf("node %q (parent: %v) removed %v: %v", node.Key(), np.Key(), removed, node)
	conf.Set("app.server.tls.priv", "/tmp/private-key.pem")
	t.Logf("tests ends.")
}

//

//

func TestStoreS_Get(t *testing.T) {
	trie := newBasicStore()
	fmt.Println(trie.MustInt("app.dump"))
	fmt.Println(trie.MustString("app.dump"))
	fmt.Println(trie.MustBool("app.dump"))
	// Output:
	// 3
	// 3
	// true
	assertEqual(t, 3, trie.MustInt("app.dump"))
	assertEqual(t, "3", trie.MustString("app.dump"))
	assertEqual(t, true, trie.MustBool("app.dump"))
}

func TestStoreS_Dump(t *testing.T) {
	conf := New()
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []any{"a", 1, false})

	data, found := conf.Get("app.logging.rotate")
	assertEqual(t, 6, data)
	assertEqual(t, true, found)
	data = conf.MustGet("app.logging.rotate")
	assertEqual(t, 6, data)

	assertEqual(t, 3, conf.MustInt("app.dump"))
	assertEqual(t, "3", conf.MustString("app.dump"))
	assertEqual(t, true, conf.MustBool("app.dump"))

	s2, e2 := conf.GetStringSlice("app.logging.words")
	assertEqual(t, nil, e2)
	assertEqual(t, []string{"a", "1", "false"}, s2)
	assertEqual(t, []string{"a", "1", "false"}, conf.MustStringSlice("app.logging.words"))
	assertEqual(t, []int{0, 1, 0}, conf.MustIntSlice("app.logging.words"))
	assertEqual(t, map[string]any{"words": []any{"a", 1, false}}, conf.MustM("app.logging.words"))

	t.Logf("%v", conf.Dump())
}

//

func assertEqual(t testing.TB, expect, actual any, msg ...any) { //nolint:govet //it's a printf/println dual interface
	if reflect.DeepEqual(expect, actual) {
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

	t.Fatalf("assertEqual failed: %v\n    expect: %v\n    actual: %v\n", mesg, expect, actual)
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

func newBasicStore(opts ...Opt) *storeS {
	conf := New(opts...)
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
	return conf.(*storeS)
}
