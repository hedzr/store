package store

import (
	"bytes"
	"sync"

	"github.com/hedzr/store/internal/radix"
)

func newStore(opts ...Opt) *storeS {
	s := &storeS{
		Trie: radix.NewTrie[any](),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithDelimiter(delimiter rune) Opt {
	return func(s *storeS) {
		s.SetDelimiter(delimiter)
	}
}

func WithPrefix(prefix string) Opt {
	return func(s *storeS) {
		s.SetPrefix(prefix)
	}
}

func WithOnChangeHandlers(handlers ...OnChangeHandler) Opt {
	return func(s *storeS) {
		s.onChangeHandlers = append(s.onChangeHandlers, handlers...)
	}
}

func WithOnNewHandlers(handlers ...OnNewHandler) Opt {
	return func(s *storeS) {
		s.onNewHandlers = append(s.onNewHandlers, handlers...)
	}
}

func WithOnDeleteHandlers(handlers ...OnDeleteHandler) Opt {
	return func(s *storeS) {
		s.OnDeleteHandlers = append(s.OnDeleteHandlers, handlers...)
	}
}

func WithFlattenSlice(b bool) Opt {
	return func(s *storeS) {
		s.flattenSlice = b
	}
}

type Opt func(s *storeS) // Opt(ions) for New Store

type Peripheral interface {
	Close()
}

// storeS is a in-memory key-value container with tree structure.
// The keys are typically dotted to represent the tree position.
type storeS struct {
	radix.Trie[any]
	loading          int32
	closers          []Peripheral
	onChangeHandlers []OnChangeHandler
	onNewHandlers    []OnNewHandler
	OnDeleteHandlers []OnDeleteHandler
	flattenSlice     bool
	parent           *storeS
}

// OnChangeHandler is called back when user setting key & value.
//
// mergingMapOrLoading is true means that user is setting key
// recursively with a map (via [Store.Merge]), or a loader
// (re-)loading its source.
type OnChangeHandler func(path string, value, oldValue any, mergingMapOrLoading bool)
type OnNewHandler func(path string, value any, mergingMapOrLoading bool)
type OnDeleteHandler func(path string, value any, mergingMapOrLoading bool)

const initialPrefixBufferSize = 64

var prefixJointPool = sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, initialPrefixBufferSize))
}}

func (s *storeS) poolGet() *bytes.Buffer {
	return prefixJointPool.Get().(*bytes.Buffer)
}

func (s *storeS) deferPoolGet(bb *bytes.Buffer) {
	bb.Reset()
	prefixJointPool.Put(bb)
}

func (s *storeS) join(args ...string) (ret string) {
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
			bb.WriteByte(byte(s.Delimiter()))
		}
		if it != "" {
			bb.WriteString(it)
			i++
		}
	}
	return bb.String()
}

// Close cleanup the internal resources.
// See [basics.Peripheral] for more information.
func (s *storeS) Close() {
	for _, c := range s.closers {
		c.Close()
	}
}

// MustGet is a shortcut to Get without error returning.
func (s *storeS) MustGet(path string) (data any) {
	var branch, found bool
	var err error
	data, branch, found, err = s.Trie.Query(path)
	if !found {
		if err != nil || !branch {
			data = nil
		}
	}
	return
}

// Get the value at path point 'path'.
func (s *storeS) Get(path string) (data any, found bool) {
	data, _, found, _ = s.Trie.Query(path)
	return
}

// Set sets key('path') and value pair into storeS.
func (s *storeS) Set(path string, data any) (oldData any) {
	old, branch, found, err := s.Trie.Query(path)
	if !found {
		if err != nil || !branch {
			old = nil
		}
	}
	if old != nil {
		oldData = old
	}

	oldData = s.setKV(path, data, !found)
	// s.tryOnSet(path, false, old, data)
	return
}

// Merge a map at path point 'pathAt'.
func (s *storeS) Merge(pathAt string, data map[string]any) (err error) {
	_, _, _, err = s.Trie.Query(pathAt)
	// if !found {
	// 	if err1 != nil || !branch {
	// 		old = nil
	// 	}
	// }
	if err != nil {
		return
	}

	err = s.loadMap(data, pathAt, false)
	// s.tryOnSet(pathAt, true, old, data)
	return
}

func (s *storeS) setKV(path string, data any, createOrModify bool) (oldData any) {
	oldData = s.Trie.Insert(path, data)
	loading := s.inLoading()
	s.tryOnSet(path, !loading, oldData, data, createOrModify)
	return
}

func (s *storeS) tryOnSet(path string, user bool, oldData, data any, createOrModify bool) {
	ptr := s

	if createOrModify {
	retryPN:
		for _, cb := range ptr.onNewHandlers {
			if cb != nil {
				cb(path, data, user)
			}
		}
		if ptr.parent != nil {
			ptr = ptr.parent
			goto retryPN
		}
		return
	}

retryPM:
	for _, cb := range ptr.onChangeHandlers {
		if cb != nil {
			cb(path, data, oldData, user)
		}
	}
	if ptr.parent != nil {
		ptr = ptr.parent
		goto retryPM
	}
}

func (s *storeS) tryOnDelete(path string, user bool, oldData any) {
	ptr := s
retryPD:
	for _, cb := range ptr.OnDeleteHandlers {
		if cb != nil {
			cb(path, oldData, user)
		}
	}
	if ptr.parent != nil {
		ptr = ptr.parent
		goto retryPD
	}
}

func (s *storeS) Remove(path string) (removed bool) {
	var rmn radix.Node[any]
	rmn, removed = s.Trie.RemoveEx(path)
	if removed {
		loading := s.inLoading()
		data := rmn.Data()
		s.tryOnDelete(path, !loading, data)
	}
	return
}

// Has tests if the given path exists
func (s *storeS) Has(path string) (found bool) {
	return s.Trie.Search(path)
}

// Locate provides an advanced interface for locating a path.
func (s *storeS) Locate(path string) (node radix.Node[any], branch, partialMatched, found bool) {
	return s.Trie.Locate(path)
}

// Dump prints internal data tree for debugging
func (s *storeS) Dump() (text string) {
	return s.Trie.Dump()
}

func (s *storeS) Clone() (newStore *storeS) { return s.Dup() } // make a clone for this store

// Dup is a native Clone tool.
//
// After Dup, a copy of original store will be created, but closers not.
// Most of the closers are cleanup code fragments coming
// from Load(WithProvider()), some of them needs to shut down the
// remote connection such as what want to do by consul provider.
//
// At this scene, the parent store still holds the cleanup closers.
func (s *storeS) Dup() (newStore *storeS) {
	newStore = &storeS{Trie: s.Trie.Dup(), flattenSlice: s.flattenSlice}
	return
}

// WithPrefix makes a lightweight copy from current storeS.
//
// The new copy is enough lite so that you can always use it with
// quite a low price.
//
// WithPrefix appends an extra prefix at end of current prefix.
// For example, on a store with old prefix "app",
// WithPrefix("store") will return a new store 'NS' with prefix
// "app.server". And NS.MustGet("type") retrieve value at key path
// "app.server.type" now.
//
//	conf := store.New()
//	s1 := conf.WithPrefix("app")
//	ns := s1.WithPrefix("server")
//	println(ns.MustGet("type"))     # print conf["app.server.type"]
//
// It simplify biz-logic codes sometimes.
//
// A [Delimiter] will be inserted at jointing prefix and key. Also at
// jointing old and new prefix.
func (s *storeS) WithPrefix(prefix string) (newStore *storeS) {
	return &storeS{parent: s, Trie: s.Trie.WithPrefix(prefix), flattenSlice: s.flattenSlice}
	// return s.withPrefixR(prefix)
}

// WithPrefixReplaced is similar with WithPrefix but it replace old
// prefix with new one instead of appending it.
//
//	conf := store.New()
//	s1 := conf.WithPrefix("app")
//	ns := s1.WithPrefixReplaced("app.server")
//	println(ns.MustGet("type"))     # print conf["app.server.type"]
//
// A [Delimiter] will be inserted at jointing prefix and key.
func (s *storeS) WithPrefixReplaced(prefix string) (newStore *storeS) {
	return &storeS{parent: s, Trie: s.Trie.WithPrefixReplaced(prefix), flattenSlice: s.flattenSlice}
}

// SetPrefix updates the prefix in current storeS.
func (s *storeS) SetPrefix(prefix string) {
	s.Trie.SetPrefix(prefix)
}
