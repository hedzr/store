package radix

import (
	"bytes"
	"sync"

	logz "github.com/hedzr/logg/slog"

	"gopkg.in/hedzr/errors.v3"
)

func NewTrie[T any]() *trieS[T] {
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

func newTrie[T any]() *trieS[T] {
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

type trieS[T any] struct {
	root      *nodeS[T]
	prefix    string
	delimiter rune
}

const dotChar rune = '.'
const initialPrefixBufferSize = 64
const maxPrefixBufferSize = 64*1024*1024 - initialPrefixBufferSize

var prefixJointPool = sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, initialPrefixBufferSize))
}}

func (s *trieS[T]) poolGet() *bytes.Buffer {
	return prefixJointPool.Get().(*bytes.Buffer)
}

func (s *trieS[T]) deferPoolGet(bb *bytes.Buffer) {
	bb.Reset()
	if bb.Cap() < maxPrefixBufferSize {
		prefixJointPool.Put(bb)
	}
}

func (s *trieS[T]) join(args ...string) (ret string) {
	switch len(args) {
	case 0:
		return
	case 1:
		return args[0]
	}

	if args[0] == "" {
		return s.join(args[1:]...)
	}

	i, bb := 0, s.poolGet()
	defer s.deferPoolGet(bb)

	for _, it := range args {
		if i > 0 {
			bb.WriteByte(byte(s.delimiter))
		}
		if it != "" {
			bb.WriteString(it)
			i++
		}
	}
	return bb.String()
}

func (s *trieS[T]) Insert(path string, data T) (oldData any) { return s.Set(path, data) }
func (s *trieS[T]) Set(path string, data T) (oldData any) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	return s.root.insert([]rune(path), path, data)
}

func (s *trieS[T]) SetComment(path, description, comment string) (ok bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched := s.search(path)
	if ok = node != nil || partialMatched; ok {
		node.description, node.comment = description, comment
	}
	return
}

func (s *trieS[T]) SetTags(path string, tags any) (ok bool) { // set extra notable data bound to a key
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched := s.search(path)
	if ok = node != nil || partialMatched; ok {
		node.tags = tags
	}
	return
}

// Merge a map at path point 'pathAt'
func (s *trieS[T]) Merge(pathAt string, data map[string]any) (err error) {
	if s.prefix != "" {
		pathAt = s.join(s.prefix, pathAt)
	}
	err = s.withPrefixR(pathAt).loadMap(data)
	return
}

func (s *trieS[T]) StartsWith(path string) (yes bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched := s.search(path)
	yes = node != nil || partialMatched
	return
}

func (s *trieS[T]) Search(path string) (found bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched := s.search(path)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

func (s *trieS[T]) Locate(path string) (node *nodeS[T], branch, partialMatched, found bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched = s.search(path)
	found, branch = node != nil && !partialMatched, safeIsBranch(node)
	return
}

func safeIsBranch[T any](node *nodeS[T]) bool { return node != nil && node.isBranch() }

func (s *trieS[T]) Has(path string) (found bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, _, partialMatched := s.search(path)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

func (s *trieS[T]) Remove(path string) (removed bool) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
	}
	node, parent, partialMatched := s.search(path)
	found := node != nil && !partialMatched // && !node.isBranch()
	if found {
		if parent != nil {
			removed = parent.remove(node)
		} else {
			logz.Warn("if given path found and return node, its parent MUST NOT be nil", "node", node, "parent", parent)
		}
	}
	return
}

func (s *trieS[T]) MustGet(path string) (data T) {
	var branch, found bool
	data, branch, found, _ = s.Query(path)
	if !found && !branch {
		data = *new(T)
	}
	return
}

func (s *trieS[T]) Get(path string) (data T, found bool) {
	data, _, found, _ = s.Query(path)
	return
}

func (s *trieS[T]) Query(path string) (data T, branch, found bool, err error) {
	if s.prefix != "" {
		path = s.join(s.prefix, path)
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
	if !found {
		err = errors.NotFound
	}
	return
}

func (s *trieS[T]) search(word string) (found, parent *nodeS[T], partialMatched bool) {
	found = s.root
	if matched, pm, child, p := found.matchR([]rune(word), s.delimiter, nil); matched || pm {
		return child, p, pm
	}
	return
}

func (s *trieS[T]) Delimiter() rune { return s.delimiter }
func (s *trieS[T]) SetDelimiter(delimiter rune) {
	s.delimiter = delimiter
}

func (s *trieS[T]) endsWith(str string, ch rune) bool {
	if str != "" {
		runes := []rune(str)
		return runes[len(runes)-1] == ch
	}
	return false
}

func (s *trieS[T]) Dump() string             { return s.root.dump(false) }
func (s *trieS[T]) dump(noColor bool) string { return s.root.dump(noColor) }

func (s *trieS[T]) Dup() (newTrie *trieS[T]) {
	newTrie = &trieS[T]{root: s.root.Dup(), prefix: s.prefix, delimiter: s.delimiter}
	return
}
