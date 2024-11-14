package radix

import (
	"testing"
)

func TestNodeS_Tag(t *testing.T) {
	trie := newBasicStore()
	trie.SetComment("app.dump", "desc", "comment")
	trie.SetTag("app.dump", 3.13)
	t.Logf("\nPath\n%v\n", trie.Dump())

	node, branch, pm, found := trie.Locate("app.dump", nil)
	if !found {
		t.Fatalf("not found: partial matched = %v, found = %v, branch = %v", pm, found, branch)
	}

	t.Logf("%v | %v | %v | %v", node.Modified(),
		node.Description(), node.Comment(),
		node.Tag(),
	)

	assert(node.Modified())
	node.SetModified(false)
	assert(!node.Modified())
	node.SetModified(true)
	assert(node.Modified())
	node.ToggleModified()
	assert(!node.Modified())
	node.ToggleModified()
	assert(node.Modified())

	if node.endsWith('.') {
		t.Fail()
	}
	if node.endsWithLite('.') {
		t.Fail()
	}

	assert(!node.IsLeaf())
	assert(!node.HasData())
	node.SetData(1)
	node.SetComment("desc", "comment")
	node.SetTag(2)
	assertEqual(t, 1, node.Data())
	assertEqual(t, 2, node.Tag())
	assertEqual(t, "desc", node.Description())
	assertEqual(t, "comment", node.Comment())
}

func TestNodeS_remove(t *testing.T) {
	trie := newBasicStore()
	trie.SetComment("app.dump", "desc", "comment")
	trie.SetTag("app.dump", 3.13)

	conf := trie.Dup()
	trie.Remove("app.debug")
	trie.Remove("app.logging.file")

	t.Logf("\nPath of 'trie'\n%v\n", trie.Dump())

	t.Logf("\nPath of 'conf'\n%v\n", conf.Dump())
}

func TestNodeS_SetModified(t *testing.T) {
	for i, c := range []struct {
		node   *nodeS[any]
		set    bool
		expect bool
	}{
		{&nodeS[any]{nType: 0}, true, true},
		{&nodeS[any]{nType: NTModified}, true, true},
		{&nodeS[any]{nType: 0}, false, false},
		{&nodeS[any]{nType: NTModified}, false, false},
	} {
		c.node.SetModified(c.set)
		actual := c.node.Modified()
		if actual != c.expect {
			t.Fatalf("%5d. expecting SetModified(%v) a node to return a %v, but got %v",
				i, c.set, c.expect, actual)
		}
	}
}

func TestNodeS_ToggleModified(t *testing.T) {
	for i, c := range []struct {
		node   *nodeS[any]
		expect bool
	}{
		{&nodeS[any]{nType: 0}, true},
		{&nodeS[any]{nType: NTModified}, false},
	} {
		c.node.ToggleModified()
		actual := c.node.Modified()
		if actual != c.expect {
			t.Fatalf("%5d. expecting ToggleModified a node to return a %v, but got %v",
				i, c.expect, actual)
		}
	}
}

func TestNodeS_ResetModified(t *testing.T) {
	for i, c := range []struct {
		node   *nodeS[any]
		expect bool
	}{
		{&nodeS[any]{nType: 0}, false},
		{&nodeS[any]{nType: NTModified}, false},
	} {
		c.node.ResetModified()
		actual := c.node.Modified()
		if actual != c.expect {
			t.Fatalf("%5d. expecting ResetModified a node to return a %v, but got %v",
				i, c.expect, actual)
		}
	}
}
