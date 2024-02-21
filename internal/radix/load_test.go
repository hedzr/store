package radix

import (
	"reflect"
	"testing"
)

func TestTrieS_Prefix(t *testing.T) {
	trie := newTrieTree2()

	trie.Insert("/about-us/legal/privacy", "what you want?")

	trie.SetDelimiter('/')
	t.Logf("\nPath of 'trie' (delimeter=%v)\n%v\n",
		trie.Delimiter(),
		trie.Dump())

	conf := trie.WithPrefix("/about-us")
	t.Logf("\nPath of 'conf' (delimeter=%v, prefix=%v)\n%v\n",
		conf.Delimiter(),
		conf.Prefix(),
		conf.Dump())

	conf.SetComment("team", "desc1", "comment1")
	conf.SetTag("team", 3.1313)
	i := conf.MustInt("team")
	if !reflect.DeepEqual(i, 4) {
		t.Fail()
	}

	n, b, pm, found := conf.Locate("team")
	if !found {
		t.Fail()
		_, _, _ = n, b, pm
	}

	t.Logf("tag of 'team': %v", n.Tag())
	if n.Tag() != 3.1313 {
		t.Fail()
	}

	// For the WithPrefix and WithPrefixReplaced:

	cfg1 := conf.WithPrefixReplaced("/about-us/legal")
	s := cfg1.MustString("privacy")
	t.Logf("'privacy': %v", s)
	if s != "what you want?" {
		t.Fail()
	}

	cfg2 := conf.WithPrefix("legal")
	s = cfg2.MustString("privacy")
	t.Logf("'privacy': %v", s)
	if s != "what you want?" {
		t.Fail()
	}

	cfg1.SetPrefix("/about-us")
	s = cfg1.MustString("legal/privacy")
	t.Logf("'privacy': %v", s)
	if s != "what you want?" {
		t.Fail()
	}
}
