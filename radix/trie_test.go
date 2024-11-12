package radix

import (
	"reflect"
	"testing"
)

func TestStoreS_Join(t *testing.T) {
	// join ...

	conf := newTrieTree()
	assertEqual(t, "", conf.Join())
	assertEqual(t, "", conf.Join(""))
	assertEqual(t, "x", conf.Join("x", ""))
	assertEqual(t, "x.y", conf.Join("x", "y"))
	assertEqual(t, "x.y", conf.Join("x", "y", ""))
	assertEqual(t, "x.y.z", conf.Join("x", "y", "z"))
	assertEqual(t, "x.y.z", conf.Join("x", "y", "z", ""))
	assertEqual(t, "x.y.z", conf.Join("", "x", "y", "z", ""))
}

func TestStoreS_join1(t *testing.T) {
	// join ...

	trie := newTrieTree()

	conf := trie.withPrefixImpl("")
	assertEqual(t, "", conf.Join())
	assertEqual(t, "", conf.Join(""))
	assertEqual(t, "x", conf.Join("x", ""))
	assertEqual(t, "x.y", conf.Join("x", "y"))
	assertEqual(t, "x.y", conf.Join("x", "y", ""))
	assertEqual(t, "x.y.z", conf.Join("x", "y", "z"))
	assertEqual(t, "x.y.z", conf.Join("x", "y", "z", ""))
	assertEqual(t, "x.y.z", conf.Join("", "x", "y", "z", ""))

	conf = trie.withPrefixImpl("A")
	assertEqual(t, "A", conf.join1("A"))
	assertEqual(t, "A", conf.join1("A", ""))
	assertEqual(t, "A.x", conf.join1("A", "x", ""))
	assertEqual(t, "A.x.y", conf.join1("A", "x", "y"))
	assertEqual(t, "A.x.y", conf.join1("A", "x", "y", ""))
	assertEqual(t, "A.x.y.z", conf.join1("A", "x", "y", "z"))
	assertEqual(t, "A.x.y.z", conf.join1("A", "x", "y", "z", ""))
	assertEqual(t, "A.x.y.z", conf.join1("A", "", "x", "y", "z", ""))
}

func TestTrieS_InsertPath(t *testing.T) {
	StatesEnvSetColorMode(true) // to disable ansi escape sequences in dump output

	trie := newTrieTree2()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	expect := `  /                             <B>
    s                           <B>
      earch/*kwargs             <L> /search/*kwargs => 1
      upport                    <L> /support => 2
    blog/:year/:month/:post     <L> /blog/:year/:month/:post => 3
    about-us/                   <B>
      team                      <L> /about-us/team => 4 // comment | tag = 3.13 ~ desc
      legal                     <L> /about-us/legal => 6
    contact                     <L> /contact => 5
`
	if ret != expect {
		t.Fatalf("failed insertions")
	}
}

func TestTrieS_Delimiter(t *testing.T) { //nolint:revive
	trie := newTrieTree2()

	// conf := trie.Dup()
	// trie.Remove("app.debug")
	// trie.Remove("app.logging.file")

	trie.SetDelimiter('/')
	t.Logf("\nPath of 'trie' (delimeter=%v)\n%v\n",
		trie.Delimiter(),
		trie.Dump())

	data, err := trie.GetM("/about-us")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	t.Logf("got: %v", data)

	if !reflect.DeepEqual(data, map[string]any{"legal": 6, "team": 4}) {
		t.Fail()
	}

	assertEqual(t, true, trie.Has("/"))
	assertEqual(t, true, trie.HasPart("/sup"))
	assertEqual(t, true, trie.Has("/support"))
	assertEqual(t, false, trie.Has("/suspend"))

	conf := trie.WithPrefix("/about-us")
	conf.SetComment("team", "desc1", "comment1")
	conf.SetTag("team", 3.1313)
	i := conf.MustInt("team")
	if !reflect.DeepEqual(i, 4) {
		t.Fail()
	}

	_ = conf.Merge("y", map[string]any{"tr": 1})

	t.Logf("\nPath of 'trie' (delimeter=%v)\n%v\n",
		trie.Delimiter(),
		trie.Dump())

	assertEqual(t, true, conf.HasPart("y"))
	assertEqual(t, true, conf.HasPart("y/tr"))
	assertEqual(t, false, conf.HasPart("y/tr/1"))

	d := conf.MustGet("x")
	assertEqual(t, nil, d)

	ia, _ := conf.Get("team")
	if !reflect.DeepEqual(ia, 4) {
		t.Fail()
	}
	ia = conf.MustGet("team")
	if !reflect.DeepEqual(ia, 4) {
		t.Fail()
	}

	conf.Remove("legal")

	has := conf.Has("legal")
	if has {
		t.Fail()
	}

	has = conf.Search("legal")
	if has {
		t.Fail()
	}
	n, _, b, pm, found := conf.Locate("team")
	if !found {
		t.Fail()
		_, _, _ = n, b, pm
	}

	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)

	// SetTag on an unexisted node must return false.
	ok := trie.SetTag("team", 3.1313)
	if ok {
		t.Fail()
	}

	// t.Logf("\nPath of 'conf'\n%v\n", conf.Dump())

	trie.Walk("", func(path, fragment string, node Node[any]) {
		t.Logf("%v %v | %v", path, fragment, node)
	})

	trie.Walk("/about-us", func(path, fragment string, node Node[any]) {
		t.Logf("%v %v | %v", path, fragment, node)
	})
}

func TestTrieS_Locate(t *testing.T) {
	trie := newTrieTree2()
	trie.SetDelimiter('/')
	t.Logf("\nPath of 'trie' (delimeter=%v)\n%v\n",
		trie.Delimiter(),
		trie.Dump())

	node, kvp, br, pm, found := trie.Locate("/search/any/thing/here")
	if !found {
		t.Fail()
	} else if br || pm {
		t.Fail()
	} else {
		t.Logf("node matched ok: %v", node)
		if kvp == nil {
			t.Fail()
		}
		if kvp["kwargs"] != "any/thing/here" {
			t.Fail()
		}
	}

	node, kvp, br, pm, found = trie.Locate("/blog/2011/09/why-so-concise")
	if !found {
		t.Fail()
	} else if br || pm {
		t.Fail()
	} else {
		t.Logf("node matched ok: %v", node)
		if kvp == nil {
			t.Fail()
		}
		if kvp["year"] != "2011" {
			t.Fail()
		}
		if kvp["month"] != "09" {
			t.Fail()
		}
		if kvp["post"] != "why-so-concise" {
			t.Fail()
		}
	}
}

func newTrieTree2() *trieS[any] {
	trie := NewTrie[any]()
	trie.Insert("/search/*kwargs", 1)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("/support", 2)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("/blog/:year/:month/:post", 3)
	trie.Insert("/about-us/team", 4)
	trie.Insert("/contact", 5)
	trie.Insert("/about-us/legal", 6)

	trie.SetComment("/about-us/team", "desc", "comment")
	trie.SetTag("/about-us/team", 3.13)

	return trie
}
