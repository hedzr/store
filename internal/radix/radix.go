package radix

import (
	"github.com/hedzr/store/ctx"
)

// Trie tree, radix-tree
type Trie[T any] interface {
	// Insert and Search, more Basic Trie Operations

	Insert(path string, data T) (oldData any)                                // Insert data (T) to path
	StartsWith(word string) (yes bool)                                       // tests if word exists, even partial matched
	Search(word string) (found bool)                                         // tests if word exists (= Has)
	Query(path string) (data T, branch, found bool, err error)               // full ability word searching (=enhanced Has)
	Locate(path string) (node *nodeS[T], branch, partialMatched, found bool) // Locate is an enhanced Has and returns more internal information (=enhanced Has)
	SetComment(path, description, comment string) (ok bool)                  // set extra meta-info bound to a key
	SetTag(path string, tags any) (ok bool)                                  // set extra notable data bound to a key
	Dump() string                                                            // dumping the node tree for debugging, including some internal states

	// Remove and Merge, Special Operations for storeS

	Remove(path string) (removed bool)                                // Remove a key and its children
	RemoveEx(path string) (nodeRemoved, parent Node[T], removed bool) // RemoveEx a key and its children

	Merge(pathAt string, data map[string]any) (err error) // advanced operation to Merge hierarchical data

	Set(path string, data T) (node Node[T], oldData any) // = Insert
	Has(path string) (found bool)                        // = Search
	Get(path string) (data T, found bool)                // shortcut to Query
	MustGet(path string) (data T)                        // shortcut to Get

	TypedGetters[T] // getters

	WithPrefix(prefix string) (entry Trie[T])         // appends prefix string and make a new instance of Trie[T]
	WithPrefixReplaced(prefix string) (entry Trie[T]) // make a new instance of Trie with prefix
	SetPrefix(prefix string)                          // set prefix. Change it on a store takes your own advantages.
	Prefix() string                                   // return current prefix string
	Delimiter() rune                                  // return current delimiter, generally it's dot ('.')
	SetDelimiter(delimiter rune)                      // setter. Change it in runtime doesn't update old delimiter inside tree nodes.

	Dup() (newTrie *trieS[T]) // a native Clone function

	Walk(path string, cb func(prefix, key string, node Node[T]))
}

type Node[T any] interface {
	isBranch() bool
	hasData() bool
	endsWith(ch rune) bool
	endsWithLite(ch rune) bool
	insert(word []rune, fullPath string, data T) (node Node[T], oldData any)
	remove(item *nodeS[T]) (removed bool)
	matchR(word []rune, delimiter rune, parentNode *nodeS[T]) (matched, partialMatched bool, child, parent *nodeS[T])
	dump(noColor bool) string

	Walk(cb func(prefix, key string, node Node[T]))

	Dup() (newNode *nodeS[T])

	Data() T

	Key() string
	Description() string
	Comment() string
	Tag() any

	Modified() bool     // node data changed by user?
	SetModified(b bool) // set modified state
	ToggleModified()    // toggle modified state
}

const NoDelimiter rune = 0

type HandlersChain func(c ctx.Ctx, next Handler)

type Handler func(c ctx.Ctx)
