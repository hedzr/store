package radix

import (
	"bytes"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"
)

// NewTrie returns a Trie-tree instance.
func NewTrie[T any]() *trieS[T] {
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

var _ Trie[any] = (*trieS[any])(nil) // assertion helper

func newTrie[T any]() *trieS[T] { //nolint:revive
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

type trieS[T any] struct {
	root      *nodeS[T]
	prefix    string
	delimiter rune
}

func (s *trieS[T]) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("Trie{\"")
	_, _ = sb.WriteString(s.prefix)
	_, _ = sb.WriteString("\", delimiter:'")
	_, _ = sb.WriteRune(s.delimiter)
	_, _ = sb.WriteString("'}")
	return sb.String()
}

func (s *trieS[T]) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	_, _ = sb.WriteString("\"trie\":{\"prefix\":")
	_, _ = sb.WriteString(strconv.Quote(s.prefix))
	_, _ = sb.WriteString(",\"delimiter\":\"")
	_, _ = sb.WriteRune(s.delimiter)
	_, _ = sb.WriteString("\"}")
	return []byte(sb.String()), nil
}

func (s *trieS[T]) dupS(root *nodeS[T], prefix string) (newTrie *trieS[T]) { //nolint:revive
	newTrie = &trieS[T]{root: root, prefix: prefix, delimiter: s.delimiter}
	return
}

//

//

//

const (
	dotChar                 rune = '.'
	initialPrefixBufferSize      = 64
	maxPrefixBufferSize          = 64*1024*1024 - initialPrefixBufferSize
)

var prefixJointPool = sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, initialPrefixBufferSize))
}}

func (s *trieS[T]) poolGet() *bytes.Buffer {
	return prefixJointPool.Get().(*bytes.Buffer) //nolint:revive
}

func (s *trieS[T]) deferPoolGet(bb *bytes.Buffer) {
	bb.Reset()
	if bb.Cap() < maxPrefixBufferSize {
		prefixJointPool.Put(bb)
	}
}

func (s *trieS[T]) UsePool(fn func(bb *bytes.Buffer, bytes int)) string {
	i, bb := 0, s.poolGet()
	defer s.deferPoolGet(bb)
	fn(bb, i)
	return bb.String()
}

func (s *trieS[T]) join1(pre string, args ...string) (ret string) {
	if pre == "" {
		return s.Join(args...)
	}

	if len(args) == 0 {
		return pre
	}

	i, bb := 0, s.poolGet()
	defer s.deferPoolGet(bb)

	_, _ = bb.WriteString(pre)

	for _, it := range args {
		if it != "" {
			_ = bb.WriteByte(byte(s.delimiter))
			_, _ = bb.WriteString(it)
			i++
		}
	}
	return bb.String()
}

func (s *trieS[T]) Join(args ...string) (ret string) {
	switch len(args) {
	case 0:
		return
	case 1:
		return args[0]
	}

	if args[0] == "" {
		return s.Join(args[1:]...)
	}

	return s.UsePool(func(bb *bytes.Buffer, bytes int) {
		for _, it := range args {
			if it != "" {
				if bytes > 0 {
					_ = bb.WriteByte(byte(s.delimiter))
				}
				_, _ = bb.WriteString(it)
				bytes++ //nolint:revive
			}
		}
	})

	// i, bb := 0, s.poolGet()
	// defer s.deferPoolGet(bb)
	//
	// for _, it := range args {
	// 	if it != "" {
	// 		if i > 0 {
	// 			_ = bb.WriteByte(byte(s.delimiter))
	// 		}
	// 		_, _ = bb.WriteString(it)
	// 		i++
	// 	}
	// }
	// return bb.String()
}

func (s *trieS[T]) Insert(path string, data T) (oldData any) { //nolint:revive
	_, oldData = s.Set(path, data)
	return
}

// Set sets the Data field into a node specified by path.
//
// If the given path cannot be found, a new node will be created at that
// location so that the new data value can be set into it.
func (s *trieS[T]) Set(path string, data T) (node Node[T], oldData any) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	return s.root.insert([]rune(path), path, data)
}

// SetComment sets the Desc and Comment field of a node specified by path.
//
// Nothing happens if the given path cannot be found.
func (s *trieS[T]) SetComment(path, description, comment string) (ok bool) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	if ok = node != nil || partialMatched; ok {
		node.description, node.comment = description, comment
	}
	return
}

// SetTag sets the Tag field of a node specified by path.
//
// Nothing happens if the given path cannot be found.
func (s *trieS[T]) SetTag(path string, tag any) (ok bool) { //nolint:revive// set extra notable data bound to a key
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	if ok = node != nil || partialMatched; ok {
		node.tag = tag
	}
	return
}

// Merge a map at path point 'pathAt'
func (s *trieS[T]) Merge(pathAt string, data map[string]any) (err error) {
	// if s.prefix != "" {
	// 	pathAt = s.Join(s.prefix, pathAt) //nolint:revive
	// }
	err = s.withPrefixImpl(pathAt).loadMap(data)
	return
}

// StartsWith tests if a path exists.
//
// Using Location to retrieve more info for seaching a path.
func (s *trieS[T]) StartsWith(path string) (yes bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	yes = node != nil || partialMatched
	if partialMatched {
		yes = strings.HasPrefix(node.pathS, path)
	}
	return
}

// Search checks the path if it exists. = StartsWith
//
// Only fully-matched node are considered as FOUND.
// Which means, if a path matched partial part, suppose matching
// `a.bcd.e` in a tree has `a.bcd.ends` node, it will be
// matched as partial-match state. But it is non-FOUND.
// And if a tree has 'a.bcd.e.nds' node, FOUND returns.
//
// Using Location to retrieve more info for seaching a path.
func (s *trieS[T]) Search(path string) (found bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

// Locate checks a path if it exists.
func (s *trieS[T]) Locate(path string) (node *nodeS[T], branch, partialMatched, found bool) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched = s.search(path)
	found, branch = node != nil && !partialMatched, safeIsBranch(node)
	return
}

func safeIsBranch[T any](node *nodeS[T]) bool { return node != nil && node.isBranch() }

// Has tests of a path exists.
//
// Delimiter shall fully match a path (typically is dot '.').
// That means, for an existed tree:
//
//	app.lite-mode
//	app.logging
//	app.logging.enabled
//	app.logging.file
//
// Has("app.logging") or Has("app.logging.enabled") will get a true.
//
// But Has("app.logging.en") returns false.
//
// And Has("app.l") must be false.
func (s *trieS[T]) Has(path string) (found bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

// Remove deleting a path from this tree.
//
// The return boolean represents there was a node removed.
// If the path does not exist, it will return false.
func (s *trieS[T]) Remove(path string) (removed bool) { //nolint:revive
	_, _, removed = s.RemoveEx(path)
	return
}

// RemoveEx deleting a path and return more status than Remove.
func (s *trieS[T]) RemoveEx(path string) (nodeRemoved, nodeParent Node[T], removed bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, parent, partialMatched := s.search(path)
	found := node != nil && !partialMatched // && !node.isBranch()
	if found {
		if parent != nil {
			removed = parent.remove(node)
			if removed {
				nodeRemoved, nodeParent = node, parent
			}
		} else {
			logz.Warn("if given path found and return node, its parent MUST NOT be nil", "node", node, "parent", parent)
		}
	}
	return
}

// MustGet is a simple Get without a checked found state.
//
// If nothing is found, zero data returned.
func (s *trieS[T]) MustGet(path string) (data T) {
	var branch, found bool
	data, branch, found, _ = s.Query(path)
	if !found && !branch {
		data = *new(T)
	}
	return
}

// Get searches the given path and return its data field if found.
func (s *trieS[T]) Get(path string) (data T, found bool) {
	data, _, found, _ = s.Query(path)
	return
}

// Query searches a path and returns the located info: 'found' boolean flag
// identify the path found or not; 'branch' flag identify the found node
// is a branch or a leaf; for a leaf node, 'data' return its Data field.
//
// If something is wrong, 'err' might collect the reason for why. But,
// it generally is errors.NotFound (errors.Code -5).
func (s *trieS[T]) Query(path string) (data T, branch, found bool, err error) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path)
	found = node != nil && !partialMatched
	if found {
		if node.isBranch() {
			branch = true
			if !node.endsWith(s.delimiter) {
				found = false
			}
		}
		if node.hasData() {
			data = node.data
		}
	}
	// if !found {
	// 	err = errors.NotFound
	// }
	err = iif(found, error(nil), error(errors.NotFound))
	return
}

func (s *trieS[T]) search(word string) (found, parent *nodeS[T], partialMatched bool) { //nolint:revive
	found = s.root
	// stringtoslicerune needs two pass full-scanning for a string, but it have to be to do.
	if matched, pm, child, p := found.matchR([]rune(word), s.delimiter, nil); matched || pm {
		return child, p, pm
	}
	found = nil
	return
}

// Delimiter returns the current delimiter in using.
//
// A trieS-tree is a radix-tree, so child nodeS may split at any
// posistion in characters of a full path.
//
// But Query or others APIs make the path is meaningful for
// the using delimiter, like dot char ('.'). This is why
// there is 'partialMatched' in Locate returning values.
//
// The delimiter can be changed at runtime. After it's changed,
// the next Locate, Query or any others will interpret the
// path with new delimiter.
//
// Since the delimiter character is part of a full path,
// we can re-interpret its meaning dynamically without extras
// costs.
//
// And the benefits are not only replaceing the delimiter
// dynamically: splitting a path by a delimiter or joining
// the segements splitted into from path are both unnecessary, in
// a conventional designing and implementing mode.
func (s *trieS[T]) Delimiter() rune { return s.delimiter }

// SetDelimiter sets the delimiter rune.
//
// When extracting a node and its data, the delimiter character is
// the decisive factor.
//
// The Store's default delimiter is '.' (dot). But you can construct
// a Trie-tree with other char, such as path separator ('/').
//
// For example,
//
//	trie := newTrie[any]()
//	trie.Insert("/search", 1)
//	trie.Insert("/support", 2)
//	trie.Insert("/blog/:post/", 3)
//	trie.Insert("/about-us/team", 4)
//	trie.Insert("/contact", 5)
//	trie.Insert("/about-us/legal", 6)
//
//	trie.SetDelimiter('/')
//	data, err := trie.GetM("/about-us")
//	assert.True(reflect.DeepEqual(data, map[string]any{"legal": 6, "team": 4}))
//
// See also TestTrieS_Delimiter(),
func (s *trieS[T]) SetDelimiter(delimiter rune) { s.delimiter = delimiter }

func (s *trieS[T]) endsWith(str string, ch rune) bool { //nolint:revive
	if str != "" {
		runes := []rune(str)
		return runes[len(runes)-1] == ch
	}
	return false
}

// Dump collects the nodes in this tree and prints them for a
// debugging purpose.
// It returns the formatted string then you can print it to
// stdout or a file.
//
//	println(trie.Dump())
//
// The dump results is decorated with ANSI escaped sequences.
// So if you want a plain pure text, enable NoColor mode
// defined in hedzr/is/states package:
//
//	import "github.com/hedzr/is/states"
//
//	states.Env().SetNoColorMode(true)
//	println(trie.Dump())
//
// Or, [StatesEnvSetColorMode(true)] can also do that.
func (s *trieS[T]) Dump() string             { return s.root.dump(false) }   //nolint:revive
func (s *trieS[T]) dump(noColor bool) string { return s.root.dump(noColor) } //nolint:revive

// Dup or Clone makes an exact deep copy of this tree.
func (s *trieS[T]) Dup() (newTrie *trieS[T]) { //nolint:revive
	return s.dupS(s.root.Dup(), s.prefix)
}

// Walk navigates the whole tree (passing "" as 'path' param) or
// a subtree from a given path.
func (s *trieS[T]) Walk(path string, cb func(path, fragment string, node Node[T])) { //nolint:revive
	root := s.root
	if path != "" {
		node, parent, partialMatched := s.search(path)
		if !partialMatched {
			root = parent
			if runes := []rune(path); runes[len(runes)-1] == s.delimiter {
				root = node
			}
		}
	}

	root.walk(0, cb)
}
