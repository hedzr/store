package store

import (
	"context"
	stderr "errors"
	"io"

	"github.com/hedzr/store/internal/radix"

	"gopkg.in/hedzr/errors.v3"
)

// New makes a new instance of storeS and return it.
//
// A storeS is a key-value container in memory with hierarchical
// tree data. A leaf or branch node can hold data. The dotted
// path
func New(opts ...Opt) *storeS { return newStore(opts...) }

func NewStoreT[T any]() StoreT[T] {
	return radix.NewTrie[T]()
}

// type storeS interface {
// 	Get(path string) (data any, found bool)
// 	Set(path string, data any)
// 	Has(path string) (found bool)
// }

// type entryS struct {
// 	name  string
// 	value Value
// }
//
// type storeSs struct {
// 	root  *entryS
// 	rootM map[string]any
//
// 	items *itemS
// }

type itemS struct { //nolint:unused
	leaves   map[string]any
	children map[string]*direntS
}

type direntS struct { //nolint:unused
	items map[string]*itemS
}

type StoreT[T any] interface {
	MustGet(path string) (data T)
	Get(path string) (data T, found bool)
	Set(path string, data T) (node radix.Node[T], oldData any)
	Has(path string) (found bool)
}

type Store interface {
	// Close cleanup the internal resources.
	// See [basics.Peripheral] for more information.
	Close()

	// MustGet is shortcut version of Get without error return.
	MustGet(path string) (data any)

	// Get the value at path point 'path'.
	Get(path string) (data any, found bool)

	// Set sets key('path') and value pair into storeS.
	Set(path string, data any) (node radix.Node[any], oldData any)

	// Remove a key and its children
	Remove(path string) (removed bool)

	// Merge a map at path point 'pathAt'.
	Merge(pathAt string, data map[string]any) (err error)

	// Has tests if the given path exists
	Has(path string) (found bool)

	// Locate provides an advanced interface for locating a path.
	Locate(path string) (node radix.Node[any], branch, partialMatched, found bool)

	radix.TypedGetters[any] // getters

	SetComment(path, description, comment string) (ok bool) // set extra meta-info bound to a key
	SetTag(path string, tags any) (ok bool)                 // set extra notable data bound to a key

	// Dump prints internal data tree for debugging
	Dump() (text string)

	// Clone makes a clone copy for this store
	Clone() (newStore *storeS)

	// Dup is a native Clone tool.
	//
	// After Dup, a copy of original store will be created, but closers not.
	// Most of the closers are cleanup code fragments coming
	// from Load(WithProvider()), some of them needs to shut down the
	// remote connection such as what want to do by consul provider.
	//
	// At this scene, the parent store still holds the cleanup closers.
	Dup() (newStore *storeS)

	// Walk iterates the whole Store.
	//
	// Walk("") walks from top-level root node.
	// Walk("app") walks from the parent of "app" node.
	// Walk("app.") walks from the "app." node.
	Walk(path string, cb func(prefix, key string, node radix.Node[any]))

	// WithPrefix makes a lightweight copy from current storeS.
	//
	// The new copy is enough light so that you can always use
	// it with quite a low price.
	//
	// WithPrefix appends an extra prefix at the end of the current
	// prefix.
	//
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
	WithPrefix(prefix ...string) (newStore *storeS) // todo need a balance on returning *storeS or Store, for WithPrefix

	// WithPrefixReplaced is similar with WithPrefix, but it replaces old
	// prefix with new one instead of appending it.
	//
	//	conf := store.New()
	//	s1 := conf.WithPrefix("app")
	//	ns := s1.WithPrefixReplaced("app.server")
	//	println(ns.MustGet("type"))     # print conf["app.server.type"]
	//
	// A [Delimiter] will be inserted at jointing prefix and key.
	//
	// todo need a balance on returning *storeS or Store, for WithPrefixReplaced.
	WithPrefixReplaced(newPrefix ...string) (newStore *storeS)

	// SetPrefix updates the prefix in current storeS.
	SetPrefix(newPrefix ...string)

	Prefix() string              // return current prefix string
	Delimiter() rune             // return current delimiter, generally it's dot ('.')
	SetDelimiter(delimiter rune) // setter. Change it in runtime doesn't update old delimiter inside tree nodes.

	// Load loads k-v pairs from external provider(s) with specified codec decoder(s).
	//
	// For those provider which run some service at background, such
	// as watching service, ctx gives a change to shut them down
	// gracefully. So you need pass a cancellable context into it.
	//
	// Or you know nothing or you don't care the terminating security,
	// simply passing context.TODO() is okay.
	Load(ctx context.Context, opts ...LoadOpt) (wr Writeable, err error)

	// WithinLoading executes a functor with loading state.
	//
	// About the Store's loading state:
	// If it's in loading, the k-v pairs will be put into store with a clean
	// modified flag.
	WithinLoading(fn func())
}

type Dumpable interface {
	Dump() string
}

var ErrNotImplemented = stderr.New("not implemented")

// The Provider gives a minimal set of interface to identify a data source.
//
// The typical data sources are: consul, etcd, file, OS environ, ....
//
// The interfaces are split to several groups: Streamable, Reader, Read, ReadBytes and Write.
//
// A provider can implement just one of the above groups.
// At this time, the other interfaces should return ErrNotImplemented.
//
// The Streamable API includes these: Keys, Count, Has, Next, Value and MustValue.
// If you are implementing it, Keys, Value and Next are Must-Have. Because our kernel
// uses Keys to confirm the provider is Streamable, and invokes Next to
// iterate the key one by one. Once a key got, Value to get its associated value.
//
// If the dataset is not very large scale, implementing Read is recommended
// to you. Read returns hierarchical data set as a nested map[string]any
// at once. Our kernel (loader) like its simple logics.
//
// Some providers may support Watchable API.
//
// All providers should always accept Codec and Position and store them.
// When changes was been monitoring by a provider, storeS will request
// a reload action and these two Properties shall be usable.
//
// Implementing OnceProvider.Write allows the provider to support Write-back mechanism.
type Provider interface {
	Read() (m map[string]any, err error) // return ErrNotImplemented as an identifier if it wants to be skipped

	ProviderSupports
}

// OnceProvider is fit for a small-scale provider.
//
// The kv data will be all loaded into memory.
type OnceProvider interface {
	ReadBytes() (data []byte, err error) // return ErrNotImplemented as an identifier if it wants to be skipped
	Write(data []byte) (err error)       // return ErrNotImplemented as an identifier if it wants to be skipped

	ProviderSupports
}

// StreamProvider is fit for a large-scale provider and load data on-demand.
type StreamProvider interface {
	Keys() (keys []string, err error)      // return ErrNotImplemented as an identifier if it wants to be skipped
	Count() int                            // count of keys and/or key-value pairs
	Has(key string) bool                   // test if the key exists
	Next() (key string, eol bool)          // return next usable key
	Value(key string) (value any, ok bool) // return the associated value
	MustValue(key string) (value any)      // return the value, or nil for non-existence key

	ProviderSupports
}

// FallbackProvider reserved for future.
type FallbackProvider interface {
	Reader() (r Reader, err error) // return ErrNotImplemented as an identifier if it wants to be skipped

	ProviderSupports
}

type ProviderSupports interface {
	GetCodec() (codec Codec)   // return the bound codec decoder
	GetPosition() (pos string) // return a position pointed to trie node path
	WithCodec(codec Codec)
	WithPosition(pos string)
}

// Reader reserved for future purpose.
type Reader interface {
	Len() int // Len returns the number of bytes of the unread portion of the slice.
	// Size returns the original length of the underlying byte slice.
	// Size is the number of bytes available for reading via ReadAt.
	// The result is unaffected by any method calls except Reset.
	Size() int64
	// Read implements the io.Reader interface.
	Read(b []byte) (n int, err error)
	// ReadAt implements the io.ReaderAt interface.
	ReadAt(b []byte, off int64) (n int, err error)
	// ReadByte implements the io.ByteReader interface.
	ReadByte() (byte, error)
	// UnreadByte complements ReadByte in implementing the io.ByteScanner interface.
	UnreadByte() error
	// ReadRune implements the io.RuneReader interface.
	ReadRune() (ch rune, size int, err error)
	// UnreadRune complements ReadRune in implementing the io.RuneScanner interface.
	UnreadRune() error
	// Seek implements the io.Seeker interface.
	Seek(offset int64, whence int) (int64, error)
	// WriteTo implements the io.WriterTo interface.
	WriteTo(w io.Writer) (n int64, err error)
	// Reset resets the Reader to be reading from b.
	Reset(b []byte)
}

// Codec is decoder and/or encoder for text format.
//
// For example, a file can be encoded with JSON format.
// So you need a JSON codec parser here.
//
// Well-known codec parsers can be JSON, YAML, TOML, ....
type Codec interface {
	Marshal(m map[string]any) (data []byte, err error)
	Unmarshal(b []byte) (data map[string]any, err error)
}

// Writeable interface
type Writeable interface {
	Save(ctx context.Context) (err error)
}

// ErrIsNotFound checks if TypedGetters returning a NotFound error.
//
//	_, err := trie.GetFloat64("app.dump.")
//	println(store.ErrIsNotFound(err))       # should be 'true'
//
// If you don't care about these errors, use MustXXX such as [radix.Trie.MustFloat64].
func ErrIsNotFound(err error) bool { return errors.Is(err, errors.NotFound) }
