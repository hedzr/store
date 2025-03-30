package radix

import (
	"time"
)

// Trie tree, an radix-tree
type Trie[T any] interface {
	// Insert and Search, more Basic Trie Operations

	Insert(path string, data T) (oldData any)                                             // Insert data (T) to path
	Search(word string) (found bool)                                                      // tests if word exists (= Has)
	Query(path string, pair KVPair) (data T, branch, found bool, err error)               // full ability word searching (=enhanced Has)
	Locate(path string, pair KVPair) (node *nodeS[T], branch, partialMatched, found bool) // Locate is an enhanced Has and returns more internal information (=enhanced Has)
	SetComment(path, description, comment string) (ok bool)                               // set extra meta-info bound to a key
	SetTag(path string, tags any) (ok bool)                                               // set extra notable data bound to a key
	Dump() string                                                                         // dumping the node tree for debugging, including some internal states

	// Remove and Merge, Special Operations for storeS

	Remove(path string) (removed bool)                                // Remove a key and its children
	RemoveEx(path string) (nodeRemoved, parent Node[T], removed bool) // RemoveEx a key and its children

	Merge(pathAt string, data map[string]any) (err error) // advanced operation to Merge hierarchical data

	StartsWith(path string, r rune) (yes bool) // tests the last path fragment by delimiter
	EndsWith(path string, r rune) (yes bool)   // tests the last path fragment by delimiter

	// Set and Get and MustGet and Has, for Store interface

	// Set = Insert
	Set(path string, data T) (node Node[T], oldData any) // = Insert
	Has(path string) (found bool)                        // = Search
	HasPart(path string) (found bool)                    // tests if word exists, even if a partial matching.
	Get(path string) (data T, found bool)                // shortcut to Query
	MustGet(path string) (data T)                        // shortcut to Get

	// GetEx gives a way to access node fields easily.
	GetEx(path string, cb func(node Node[T], data T, branch bool, kvpair KVPair))

	// SetNode at once, advanced api here.
	SetNode(path string, data T, tag any, descriptionAndComments ...string) (ret Node[T], oldData any)
	// SetEmpty clear the Data field.
	SetEmpty(path string) (oldData any)
	// Update a node whether it existed or not.
	Update(path string, cb func(node Node[T], old any))

	// SetTTL sets a ttl timeout for a branch or a leaf node.
	//
	// At ttl arrived, the leaf node value will be cleared.
	// For a branch node, it will be dropped.
	//
	// Once you're using SetTTL, don't forget call Close().
	// For example:
	//
	//	conf := newTrieTree()
	//	defer conf.Close()
	//
	//	path := "app.verbose"
	//	conf.SetTTL(path, 200*time.Millisecond, func(ctx context.Context, s *TTL[any], nd *Node[any]) {
	//		t.Logf("%q cleared", path)
	//	})
	//
	// **[Pre-API]**
	//
	// SetTTL is a prerelease API since v1.2.5, it's mutable in the
	// several future releases recently.
	//
	// The returned `state`: 0 assumed no error.
	SetTTL(path string, ttl time.Duration, cb OnTTLRinging[T]) (state int)

	// SetTTLFast ignores the existance validation of the target node
	// since it is used as a parameter.
	SetTTLFast(node Node[T], ttl time.Duration, cb OnTTLRinging[T]) (state int)

	// SetEx is advanced version of Set.
	//
	// Using it to setup a new node at once. For example:
	//
	//	conf.SetEx("app.logging.auto-stop", true,
	//	  func(path string, oldData any, node radix.Node[any], trie radix.Trie[any]) {
	//	    conf.SetTTL(path, 30*time.Minute,
	//	      func(s *radix.TTL[any], node radix.Node[any]) {
	//	        conf.Remove(node.Key()) // erase the key with the node
	//	      })
	//	    // Or:
	//	    trie.SetTTLFast(node, 3*time.Second, nil)
	//	    // Or:
	//	    node.SetTTL(3*time.Second, trie, nil)
	//	  })
	SetEx(path string, data T, cb OnSetEx[T]) (oldData any)

	TypedGetters[T] // getters

	GetTag(path string) (tag any, err error)            // get tag field directly
	MustGetTag(path string) (tag any)                   // mustget tag field directly
	GetComment(path string) (comment string, err error) // get comment field directly
	MustGetComment(path string) (comment string)        // mustget comment field directly

	WithPrefix(prefix ...string) (entry Trie[T])            // appends prefix string and make a new instance of Trie[T]
	WithPrefixReplaced(newPrefix ...string) (entry Trie[T]) // make a new instance of Trie with prefix
	SetPrefix(newPrefix ...string)                          // set prefix. Change it on a store takes your own advantages.
	Prefix() string                                         // return current prefix string
	Delimiter() rune                                        // return current delimiter, generally it's dot ('.')
	SetDelimiter(delimiter rune)                            // setter. Change it in runtime doesn't update old delimiter inside tree nodes.

	// Dup duplicates a new instance from this one. = Clone.
	Dup() (newTrie *trieS[T]) // a native Clone function

	// Walk iterators the whole tree for each node.
	Walk(path string, cb func(path, fragment string, node Node[T]))

	String() string               // for log/slog text mode
	MarshalJSON() ([]byte, error) // for log/slog json mode
}

// Node is a Trie-tree node.
type Node[T any] interface {
	// isBranch() bool
	// hasData() bool
	// endsWith(ch rune) bool
	// endsWithLite(ch rune) bool
	// insert(word []rune, fullPath string, data T) (node Node[T], oldData any)
	// remove(item *nodeS[T]) (removed bool)
	// matchR(word []rune, delimiter rune, parentNode *nodeS[T]) (matched, partialMatched bool, child, parent *nodeS[T])
	// dump(noColor bool) string

	EndsWith(ch rune) bool
	StartsWith(ch rune) bool

	// Walk iterators the whole sub-tree from this node.
	Walk(cb func(path, fragment string, node Node[T]))

	// Dup duplicates a new instance from this one. = Clone.
	Dup() (newNode *nodeS[T])

	Data() T             // retrieve the data value, just valid for leaf node
	Key() string         // retrieve the key field (full path of the node), just valid for leaf node
	Description() string // retrieve the description field, just valid for leaf node
	Comment() string     // retrieve the remarks field, just valid for leaf node
	Tag() any            // retrieve the tag field, just valid for leaf node

	SetData(data T)                  // setter for Data field
	SetEmpty()                       // SetEmpty clear the Data field. An empty node is same with node.Empty() or ! HasData()
	SetComment(desc, comment string) // setter for Description and Comment field
	SetTag(tag any)                  // setter for Tag field

	// SetTTL sets a ttl timeout for a branch or a leaf node.
	//
	// At ttl arrived, the leaf node value will be cleared.
	// For a branch node, it will be dropped.
	//
	// Once you're using SetTTL, don't forget call Close().
	SetTTL(duration time.Duration, trie Trie[T], cb OnTTLRinging[T])

	Modified() bool     // node data changed by user?
	SetModified(b bool) // set modified state
	ToggleModified()    // toggle modified state

	IsLeaf() bool   // check if a node type is leaf
	IsBranch() bool // check if a node is branch (has children)
	HasData() bool  // check if a node has data. only leaf node can contain data field. = ! Empty() bool
	Empty() bool    // check if the node has no data. It means an empty data.

	KeyPiece() string // key piece field for this node
}

const NoDelimiter rune = 0 // reserved for an internal special tree

// type HandlersChain func(c ctx.Ctx, next Handler)
//
// type Handler func(c ctx.Ctx)
