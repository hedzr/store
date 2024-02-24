package radix

import (
	"testing"
)

func TestNodeS_Tag(t *testing.T) {
	trie := newBasicStore()
	trie.SetComment("app.dump", "desc", "comment")
	trie.SetTag("app.dump", 3.13)
	t.Logf("\nPath\n%v\n", trie.Dump())

	node, branch, pm, found := trie.Locate("app.dump")
	if !found {
		t.Fatalf("not found")
		_, _ = branch, pm
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
