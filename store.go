package store

import (
	"bytes"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hedzr/store/radix"
)

func newStore(opts ...Opt) *storeS {
	_ = os.Setenv("STORE_VERSION", Version)
	s := &storeS{
		Trie: radix.NewTrie[any](),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithDelimiter sets the delimiter char.
//
// A delimiter char is generally used for extracting the key-value
// pair via GetXXX, MustXXX, e.g., MustInt, MustStringSlice, ....
func WithDelimiter(delimiter rune) Opt {
	return func(s *storeS) {
		s.SetDelimiter(delimiter)
	}
}

// WithPrefix sets the associated prefix for the tree path.
func WithPrefix(prefix string) Opt {
	return func(s *storeS) {
		s.SetPrefix(prefix)
	}
}

// WithOnChangeHandlers allows user's handlers can be callback once a node changed.
func WithOnChangeHandlers(handlers ...OnChangeHandler) Opt {
	return func(s *storeS) {
		s.onChangeHandlers = append(s.onChangeHandlers, handlers...)
	}
}

// WithOnNewHandlers allows user's handlers can be callback if a new node has been creating.
func WithOnNewHandlers(handlers ...OnNewHandler) Opt {
	return func(s *storeS) {
		s.onNewHandlers = append(s.onNewHandlers, handlers...)
	}
}

// WithOnDeleteHandlers allows user's handlers can be callback once a node removed.
func WithOnDeleteHandlers(handlers ...OnDeleteHandler) Opt {
	return func(s *storeS) {
		s.OnDeleteHandlers = append(s.OnDeleteHandlers, handlers...)
	}
}

// WithFlattenSlice sets a bool flag to tell Store the slice value should be
// treated as node leaf. The index of the slice would be part of node path.
// For example, you're loading a slice []string{"A","B"} into node path
// "app.slice", the WithFlattenSlice(true) causes the following structure:
//
//	app.slice.0 => "A"
//	app.slice.1 => "B"
//
// Also, WithFlattenSlice makes the map values to be flattened into a tree.
func WithFlattenSlice(b bool) Opt {
	return func(s *storeS) {
		s.flattenSlice = b
	}
}

// WithWatchEnable allows watching the external source if its provider
// supports Watchable ability.
func WithWatchEnable(b bool) Opt {
	return func(s *storeS) {
		s.allowWatch = b
	}
}

type Opt func(s *storeS) // Opt(ions) for New Store

// Peripheral is closeable.
type Peripheral interface {
	Close()
}

// storeS is an in-memory key-value container with tree structure supporting.
// The keys are typically dotted to represent the tree position.
type storeS struct {
	radix.Trie[any]

	loading          int32
	saving           int32
	closers          []Peripheral
	onChangeHandlers []OnChangeHandler
	onNewHandlers    []OnNewHandler
	OnDeleteHandlers []OnDeleteHandler

	// The following members need to Dup, WithPrefix, and
	// WithPrefixReplaced.
	// See dupS()

	parent *storeS

	flattenSlice bool
	allowWatch   bool
}

func (s *storeS) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("Store{")
	_, _ = sb.WriteString(s.Trie.String())
	_, _ = sb.WriteString("}")
	return sb.String()
}

func (s *storeS) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	_, _ = sb.WriteString("{")
	if b, err := s.Trie.MarshalJSON(); err != nil {
		return nil, err
	} else {
		_, _ = sb.Write(b)
	}
	_, _ = sb.WriteString("}")
	return []byte(sb.String()), nil
}

func (s *storeS) dupS(trie radix.Trie[any]) (newStore *storeS) {
	newStore = &storeS{
		Trie:         trie,
		flattenSlice: s.flattenSlice,
		allowWatch:   s.allowWatch,
		loading:      s.loading,
		// don't dup the member 'parent' here
	}
	return
}

//

//

//

var _ radix.TypedGetters[any] = (*storeS)(nil) // assertion helper

var _ Store = (*dummyS)(nil) // assertion helper

var _ MinimalStoreT[any] = (*dummyS)(nil) // assertion helper

// OnChangeHandler is called back when user setting key & value.
//
// mergingMapOrLoading is true means that user is setting key
// recursively with a map (via [Store.Merge]), or a loader
// (re-)loading its source.
type OnChangeHandler func(path string, value, oldValue any, mergingMapOrLoading bool)
type OnNewHandler func(path string, value any, mergingMapOrLoading bool)    // when user setting a new key
type OnDeleteHandler func(path string, value any, mergingMapOrLoading bool) // when user deleting a key

func (*OnChangeHandler) GobDecode([]byte) error    { return nil }
func (OnChangeHandler) GobEncode() ([]byte, error) { return nil, nil }
func (*OnNewHandler) GobDecode([]byte) error       { return nil }
func (OnNewHandler) GobEncode() ([]byte, error)    { return nil, nil }
func (*OnDeleteHandler) GobDecode([]byte) error    { return nil }
func (OnDeleteHandler) GobEncode() ([]byte, error) { return nil, nil }

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
		if it != "" {
			if i > 0 {
				bb.WriteByte(byte(s.Delimiter()))
			}
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
	data, branch, found, err = s.Trie.Query(path, nil)
	if !found {
		if err != nil || !branch {
			data = nil
		}
	}
	return
}

// Get the value at path point 'path'.
func (s *storeS) Get(path string) (data any, found bool) {
	data, _, found, _ = s.Trie.Query(path, nil)
	return
}

// Set sets key('path') and value pair into storeS.
func (s *storeS) Set(path string, data any) (node radix.Node[any], oldData any) {
	old, branch, found, err := s.Trie.Query(path, nil)
	if !found {
		if err != nil || !branch {
			old = nil
		}
	}
	if old != nil {
		oldData = old
	}

	node, oldData = s.setKV(path, data, !found, nil)
	// s.tryOnSet(path, false, old, data)
	return
}

// SetTTL sets a ttl timeout for a branch or a leaf node.
//
// At ttl arrived, the leaf node value will be cleared.
// For a branch node, it will be dropped.
//
// Once you're using SetTTL, don't forget call Close().
// For example:
//
//	conf := newBasicStore()
//	defer conf.Close()
//
//	path := "app.verbose"
//	conf.SetTTL(path, 200*time.Millisecond, func(ctx context.Context, func(s *radix.TTL[any], nd radix.Node[any]) {
//		t.Logf("%q cleared", path)
//	})
//
// **[Pre-API]**
//
// SetTTL is a prerelease API since v1.2.5, it's mutable in the
// several future releases recently.
//
// The returned `state`: 0 assumed no error.
func (s *storeS) SetTTL(path string, ttl time.Duration, cb radix.OnTTLRinging[any]) (state int) {
	state = s.Trie.SetTTL(path, ttl, cb)
	return
}

// Merge a map at path point 'pathAt'.
func (s *storeS) Merge(pathAt string, data map[string]any) (err error) {
	// _, _, _, err = s.Trie.Query(pathAt)
	// // if !found {
	// // 	if err1 != nil || !branch {
	// // 		old = nil
	// // 	}
	// // }
	// if err != nil {
	// 	return
	// }

	err = s.loadMap(data, pathAt, false, nil)
	// s.tryOnSet(pathAt, true, old, data)
	return
}

// func (s *storeS) setKValPkg(path string, vp ValPkg, createOrModify bool) (node radix.Node[any], oldData any) {
// 	s.Trie.SetComment(path, vp.Desc, vp.Comment)
// 	s.Trie.SetTag(path, vp.Tag)
//
// 	loading := s.inLoading()
// 	user := !loading
// 	if user {
// 		if oldData != nil {
// 			createOrModify = false // set it to is-modifying instead of is-creating
// 		}
// 		if node != nil {
// 			node.SetModified(true)
// 		}
// 	}
// 	s.tryOnSet(path, user, oldData, vp.Value, createOrModify)
// 	return
// }

func (s *storeS) setKV(path string, data any, createOrModify bool, onSet lmOnSet) (node radix.Node[any], oldData any) {
	node, oldData = s.Trie.Set(path, data)
	loading := s.inLoading()
	user := !loading
	if user {
		if oldData != nil {
			createOrModify = false // set it to is-modifying instead of is-creating
		}
		if node != nil {
			if onSet != nil {
				onSet(node)
			}
			node.SetModified(true)
		}
	}
	s.tryOnSet(path, user, oldData, data, createOrModify)
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

func (s *storeS) tryOnDelete(path string, user bool, oldData any, node, np radix.Node[any]) {
	ptr := s
	_, _ = node, np
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
	var rmn, np radix.Node[any]
	rmn, np, removed = s.Trie.RemoveEx(path)
	if removed {
		loading := s.inLoading()
		data := rmn.Data()
		s.tryOnDelete(path, !loading, data, rmn, np)
	}
	return
}

func (s *storeS) RemoveEx(path string) (nodeRemoved, nodeParent radix.Node[any], removed bool) {
	nodeRemoved, nodeParent, removed = s.Trie.RemoveEx(path)
	if removed {
		loading := s.inLoading()
		data := nodeRemoved.Data()
		s.tryOnDelete(path, !loading, data, nodeRemoved, nodeParent)
	}
	return
}

// Has tests if the given path exists.
func (s *storeS) Has(path string) (found bool) {
	return s.Trie.Search(path)
}

// Locate provides an advanced interface for locating a path.
func (s *storeS) Locate(path string, kvpair radix.KVPair) (node radix.Node[any], branch, partialMatched, found bool) {
	return s.Trie.Locate(path, kvpair)
}

// Dump prints internal data tree for debugging
func (s *storeS) Dump() (text string) {
	return s.Trie.Dump()
}

func (s *storeS) Clone() (newStore Store) { return s.Dup() } // make a clone for this store

// Dup is a native Clone tool.
//
// After Dup, a copy of the original store will be created, but
// closers not.
// Most of the closers are cleanup code fragments coming
// from Load(WithProvider()), some of them needs to shut down the
// remote connection such as what want to do by consul provider.
//
// At this scene, the parent store still holds the cleanup closers.
func (s *storeS) Dup() (newStore Store) {
	return s.dupS(s.Trie.Dup())
}

// WithPrefix makes a lightweight copy from current storeS.
//
// The new copy is enough light so that you can always use it with
// quite a low price.
//
// WithPrefix appends an extra prefix at the end of the current prefix.
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
func (s *storeS) WithPrefix(prefix ...string) (newStore Store) {
	return s.dupS(s.Trie.WithPrefix(prefix...))
}

// WithPrefixReplaced is similar with WithPrefix, but it replaces
// old prefix with new one instead of appending it.
//
//	conf := store.New()
//	s1 := conf.WithPrefix("app")
//	ns := s1.WithPrefixReplaced("app.server")
//	println(ns.MustGet("type"))     # print conf["app.server.type"]
//
// A [Delimiter] will be inserted at jointing prefix and key.
func (s *storeS) WithPrefixReplaced(newPrefix ...string) (newStore Store) {
	return s.dupS(s.Trie.WithPrefixReplaced(newPrefix...))
}

// SetPrefix updates the prefix in current storeS.
func (s *storeS) SetPrefix(newPrefix ...string) {
	s.Trie.SetPrefix(newPrefix...)
}

// To finds a given path and loads the subtree into
// 'holder', typically 'holder' could be a struct.
//
// For yaml input
//
//	app:
//	  server:
//	    sites:
//	      - name: default
//	        addr: ":7999"
//	        location: ~/Downloads/w/docs
//
// The following codes can load it into sitesS struct:
//
//	var sites sitesS
//	err = store.To(store.WithPrefix("app"), "server.sites", &sites)
//
//	type sitesS struct{ Sites []siteS }
//
//	type siteS struct {
//	  Name        string
//	  Addr        string
//	  Location    string
//	}
//
// In this above case, 'store' loaded yaml and built it
// into memory, and extract 'server.sites' into 'sitesS'.
// Since 'server.sites' is a yaml array, it was loaded
// as a store entry and holds a slice value, so GetSectionFrom
// extract it to sitesS.Sites field.
//
// The optional MOpt operators could be:
//   - WithKeepPrefix
//   - WithFilter
func To[T any](s Store, path string, holder *T, opts ...radix.MOpt[any]) (err error) {
	return s.To(path, holder, opts...)
}
