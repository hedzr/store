package radix

import (
	"fmt"
	"slices"
	"testing"

	"github.com/hedzr/is/states"
)

func newTrieTree() *trieS[any] {
	trie := newTrie[any]()
	trie.Insert("app.debug", 1)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("app.verbose", 2)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("app.dump", 3)
	trie.Insert("app.logging.file", 4)
	trie.Insert("app.server.start", 5)
	trie.Insert("app.logging.rotate", 6)
	trie.Insert("app.logging.words", []any{"a", 1, false})
	// trie.Insert("app.logging.words", []string{"a", "1", "false"})
	trie.Insert("app.server.sites", 1)
	return trie
}

func newBasicStore() *trieS[any] {
	conf := newTrie[any]()
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

func TestTrieS_General(t *testing.T) {
	trie := newTrie[any]()
	trie.Insert("apple", 1)
	t.Logf("\nPath    Data\n%v\n", trie.dump(true))
	assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	trie.Insert("app", 2)
	t.Logf("\nPath    Data\n%v\n", trie.dump(true))
	assertTrue(t, trie.Search("app"), `expecting trie.Search("app") again return true`) // 返回 True
}

func TestTrieS_Insert(t *testing.T) {
	states.Env().SetNoColorMode(true)

	trie := newBasicStore()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	expect := `  app.                          <B>
    d                           <B>
      ebug                      <L> app.debug => false
      ump                       <L> app.dump => 3
    verbose                     <L> app.verbose => true
    logging.                    <B>
      file                      <L> app.logging.file => /tmp/1.log
      rotate                    <L> app.logging.rotate => 6
      words                     <L> app.logging.words => [a 1 false]
    server.start                <L> app.server.start => 5
`
	if ret != expect {
		t.Fatalf("failed insertions")
	}
}

func TestTrieS_Set(t *testing.T) {
	states.Env().SetNoColorMode(true)

	trie := newBasicStore()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)

	expect := `  app.                          <B>
    d                           <B>
      ebug                      <L> app.debug => false
      ump                       <L> app.dump => 3
    verbose                     <L> app.verbose => true
    logging.                    <B>
      file                      <L> app.logging.file => /tmp/1.log
      rotate                    <L> app.logging.rotate => 6
      words                     <L> app.logging.words => [a 1 false]
    server.start                <L> app.server.start => 5
`
	if ret != expect {
		t.Fatalf("failed insertions")
	}
}

func TestTrieS_Merge(t *testing.T) {
	// trie := NewStoreT[any]()
	trie := newBasicStore()
	if err := trie.Merge("app.debug.map", map[string]any{"tag": []string{"delve", "verbose"}, "verbose": true}); err != nil {
		t.Fatalf(`Merge failed: %v`, err)
	}
	t.Logf("\nPath\n%v\n", trie.Dump())

	assertTrue(t, slices.Equal(trie.MustGet("app.debug.map.tag").([]string), []string{"delve", "verbose"}),
		`expecting trie.Get("app.debug.map.tag") return '[delve verbose]'`)
}

func TestTrieS_WithPrefix(t *testing.T) {
	trie := newBasicStore()
	t.Logf("\nPath\n%v\n", trie.Dump())

	assertTrue(t, trie.MustGet("app.logging.rotate") == 6, `expecting trie.Get("app.logging.rotate") return 6`)
}

func TestTrieS_Query(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
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
		data, _, found, _ := trie.Query(c.query)
		if found == c.found && data == c.data {
			continue
		}
		t.Fatalf("%5d. querying %q and got (%v, %v), but expecting (%v, %v)", i, c.query, found, data, c.found, c.data)
	}
}

func TestTrieS_GetR_2(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	m, err := trie.GetR("app")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('app') returns:\n\n%v\n\n", m)
	// t.Logf("GetR('app') returns:")
	// spew.Default.Println(m)

	m, err = trie.GetR("")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('') returns:\n\n%v\n\n", m)
	// t.Logf("GetR('') returns:")
	// spew.Default.Println(m)
}

func TestTrieS_GetM_2(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	m, err := trie.GetM("")
	t.Logf("map: %v | err: %v", m, err)
}

//

//

//

func assertTrue(t testing.TB, cond bool, msg ...any) { //nolint:revive
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

func assertFalse(t testing.TB, cond bool, msg ...any) { //nolint:revive
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
