package radix

import (
	"reflect"
	"testing"

	"github.com/hedzr/is/states"
)

func TestStoreS_join(t *testing.T) {
	// join ...

	conf := newTrieTree()
	assertEqual(t, "", conf.join())
	assertEqual(t, "", conf.join(""))
	assertEqual(t, "x", conf.join("x", ""))
	assertEqual(t, "x.y", conf.join("x", "y"))
	assertEqual(t, "x.y", conf.join("x", "y", ""))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z"))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z", ""))
	assertEqual(t, "x.y.z", conf.join("", "x", "y", "z", ""))
}

func TestStoreS_join1(t *testing.T) {
	// join ...

	trie := newTrieTree()

	conf := trie.withPrefix("")
	assertEqual(t, "", conf.join())
	assertEqual(t, "", conf.join(""))
	assertEqual(t, "x", conf.join("x", ""))
	assertEqual(t, "x.y", conf.join("x", "y"))
	assertEqual(t, "x.y", conf.join("x", "y", ""))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z"))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z", ""))
	assertEqual(t, "x.y.z", conf.join("", "x", "y", "z", ""))

	conf = trie.withPrefix("A")
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
	states.Env().SetNoColorMode(true)

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
      earch                     <L> /search => 1
      upport                    <L> /support => 2
    blog/:post/                 <L> /blog/:post/ => 3
    about-us/                   <B>
      team                      <L> /about-us/team => 4 // comment | tag = 3.13 ~ desc
      legal                     <L> /about-us/legal => 6
    contact                     <L> /contact => 5
`
	if ret != expect {
		t.Fatalf("failed insertions")
	}
}

func TestTrieS_Delimiter(t *testing.T) {
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

	conf := trie.WithPrefix("/about-us")
	conf.SetComment("team", "desc1", "comment1")
	conf.SetTag("team", 3.1313)
	i := conf.MustInt("team")
	if !reflect.DeepEqual(i, 4) {
		t.Fail()
	}

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
	n, b, pm, found := conf.Locate("team")
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

func newTrieTree2() *trieS[any] {
	trie := NewTrie[any]()
	trie.Insert("/search", 1)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("/support", 2)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Insert("/blog/:post/", 3)
	trie.Insert("/about-us/team", 4)
	trie.Insert("/contact", 5)
	trie.Insert("/about-us/legal", 6)

	trie.SetComment("/about-us/team", "desc", "comment")
	trie.SetTag("/about-us/team", 3.13)

	return trie
}
