package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	logz "github.com/hedzr/logg/slog"
	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store/internal/radix"
)

type loadS struct {
	*storeS
	position string
	codec    Codec
	provider Provider
}

type LoadOpt func(*loadS)

func WithProvider(provider Provider) LoadOpt {
	return func(s *loadS) {
		s.provider = provider
	}
}

func WithCodec(codec Codec) LoadOpt {
	return func(s *loadS) {
		s.codec = codec
	}
}

func WithStorePrefix(prefix string) LoadOpt {
	return func(s *loadS) {
		s.storeS = s.storeS.WithPrefixReplaced(prefix)
	}
}

func WithPosition(position string) LoadOpt {
	return func(s *loadS) {
		s.position = position
	}
}

// WithStoreFlattenSlice can destruct slice/map as tree hierarchy
// instead of treating it as a node value.
func WithStoreFlattenSlice(b bool) LoadOpt {
	return func(s *loadS) {
		s.flattenSlice = b
	}
}

// WithKeepPrefix can construct tree nodes hierarchy with the key prefix.
//
// By default, the prefix will be stripped from a given key path.
//
// For example, with a store set a prefix 'app.server',
// `store.Put("app.server.tls", map[string]any{ "certs": "some/where.pem" }` will
// produce the tree structure like:
//
//	app.
//	  Server.
//	    tls.
//	      certs   => "some/where.pem"
//
// But if you enable keep-prefix setting, the code can only be written as:
//
//	store.Put("tls", map[string]any{ "certs": "some/where.pem" }
//
// We recommend using our default setting except that you knew what you want.
// By using the default setting, i.e. keepPrefix == false, we will strip
// the may-be-there prefix if necessary. So both "app.server.tls" and "tls"
// will work properly as you really want.
func WithKeepPrefix(b bool) radix.MOpt {
	return radix.WithKeepPrefix(b)
}

func (s *storeS) inLoading() bool { return atomic.LoadInt32(&s.loading) == 1 }

func (s *storeS) Load(ctx context.Context, opts ...LoadOpt) (err error) {
	if atomic.CompareAndSwapInt32(&s.loading, 0, 1) {
		defer func() { atomic.CompareAndSwapInt32(&s.loading, 1, 0) }()

		var loader = newLoader(s, opts...)

		var data map[string]any
		data, err = loader.tryLoad(ctx) // load dataset from source via loader
		if err != nil {
			return
		}

		// merge dataset into store
		if err = loader.loadMap(data, loader.Prefix(), true); err != nil {
			return
		}

		loader.startWatch(ctx, loader)
	}
	return
}

func newLoader(st *storeS, opts ...LoadOpt) *loadS {
	var loader = &loadS{
		storeS:   st,
		codec:    nil,
		provider: nil,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(loader)
		}
	}

	if loader.codec == nil {
		loader.codec = loader.provider.GetCodec()
	} else {
		loader.provider.WithCodec(loader.codec)
	}
	if loader.position == "" {
		loader.position = loader.provider.GetPosition()
	} else {
		loader.provider.WithPosition(loader.position)
	}

	// loader.provider.WithStorePrefix(s.prefix)
	return loader
}

func (s *loadS) tryLoad(ctx context.Context) (data map[string]any, err error) {
	var b []byte

	data, err = s.provider.Read()

	if errors.Is(err, NotImplemented) {
		switch fp := s.provider.(type) {
		case OnceProvider:
			b, err = fp.ReadBytes()
		case StreamProvider:
			err = nil
			for {
				k, eol := fp.Next()
				if eol {
					break
				}
				s.setKV(k, fp.MustValue(k), true)
			}
		}
	}
	if err != nil {
		return
	}

	if s.codec != nil {
		data, err = s.codec.Unmarshal(b)
	} else if data == nil {
		err = json.Unmarshal(b, &data)
	}
	return
}

func (s *storeS) loadMap(m map[string]any, position string, creating bool) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	for k, v := range m {
		s.loadMapByValueType(ec, m, position, k, v, creating)
	}
	return
}

func (s *storeS) loadMapByValueType(ec errors.Error, m map[string]any, position, k string, v any, creating bool) {
	switch vv := v.(type) {
	case map[string]any:
		ec.Attach(s.loadMap(vv, s.join(position, k), creating))
	case []map[string]any:
		if s.flattenSlice {
			buf := make([]byte, 0, len(k)+16)
			for i, mm := range vv {
				buf = append(buf, k...)
				buf = append(buf, byte(s.Delimiter()))
				buf = strconv.AppendInt(buf, int64(i), 10)
				ec.Attach(s.loadMap(mm, s.join(position, string(buf)), creating))
				buf = buf[:0]
			}
			break
		}
		s.WithPrefixReplaced(position).setKV(k, v, creating)
	case []any:
		if s.flattenSlice {
			buf := make([]byte, 0, len(k)+16)
			for i, mm := range vv {
				// if s.prefix != "" {
				// 	buf = append(buf, s.prefix...)
				// 	buf = append(buf, byte(s.Delimiter()))
				// }
				buf = append(buf, k...)
				buf = append(buf, byte(s.Delimiter()))
				buf = strconv.AppendInt(buf, int64(i), 10)
				s.loadMapByValueType(ec, m, position, string(buf), mm, creating)
				buf = buf[:0]
			}
			break
		}
		s.WithPrefixReplaced(position).setKV(k, v, creating)
	default:
		s.WithPrefixReplaced(position).setKV(k, v, creating)
	}
	return
}

type Watchable interface {
	Watch(ctx context.Context, cb func(event any, err error)) error

	// Close provides a closer to cleanup the peripheral gracefully
	Close()
	// basics.Peripheral
}

type Change interface {
	// Key() string
	// Val() any

	Next() (key string, val any, ok bool)

	Path() string // specially for 'file' provider

	Op() Op //
	Has(op Op) bool
	Timestamp() time.Time

	Provider() Provider
}

type Op uint32 // Op describes a set of file operations.

var opStrings = map[Op]string{
	OpCreate: "create",
	OpWrite:  "modify",
	OpRename: "rename",
	OpRemove: "remove",
	OpChmod:  "chmod",
	OpNone:   "none",
}

var opStringsRev = map[string]Op{
	"create": OpCreate,
	"new":    OpCreate,
	"modify": OpWrite,
	"write":  OpWrite,
	"rename": OpRename,
	"remove": OpRemove,
	"delete": OpRemove,
	"rm":     OpRemove,
	"chmod":  OpChmod,
	"none":   OpNone,
}

func (s *Op) UnmarshalText(text []byte) error {
	// panic("implement me")
	op, ok := opStringsRev[string(text)]
	if ok {
		*s = op
		return nil
	}
	return errors.New("bad/unknown string, can't unmarshal to Op")
}

func (s Op) MarshalText() (text []byte, err error) {
	sz, ok := opStrings[s]
	if ok {
		return []byte(strings.ToUpper(sz)), nil
	}
	return []byte(fmt.Sprintf("Op(%d)", s)), nil
}

// The operations fsnotify can trigger; see the documentation on [Watcher] for a
// full description, and check them with [Event.Has].
const (
	// OpCreate is a new pathname was created.
	OpCreate Op = 1 << iota

	// OpWrite the pathname was written to; this does *not* mean the write has finished,
	// and a write can be followed by more writes.
	OpWrite

	// OpRemove the path was removed; any watches on it will be removed. Some "remove"
	// operations may trigger a Rename if the file is actually moved (for
	// example "remove to trash" is often a rename).
	OpRemove

	// OpRename the path was renamed to something else; any watched on it will be
	// removed.
	OpRename

	// OpChmod file attributes were changed.
	//
	// It's generally not recommended to take action on this event, as it may
	// get triggered very frequently by some software. For example, Spotlight
	// indexing on macOS, anti-virus software, backup software, etc.
	OpChmod

	OpNone = 0
)

func (s Op) Marshal() []byte {
	return nil
}

func (s *storeS) startWatch(ctx context.Context, loader *loadS) {
	if loader.provider == nil {
		return
	}
	if w, ok := loader.provider.(Watchable); ok {
		if err := w.Watch(ctx, s.applyExternalChanges); err != nil {
			logz.Error("[Watcher.StartWatch.ERROR]", "err", err)
		} else {
			s.closers = append(s.closers, w)
		}
	}
}

func (s *storeS) applyExternalChanges(event any, err error) {
	if err != nil {
		logz.Error("[Watcher.ERROR]", "err", err)
		return
	}

	if fse, ok := event.(Change); ok {
		s.applyChanges(fse)
	}
}

func (s *storeS) applyChanges(ev Change) {
	// if err := s.Load(WithProvider(ev.Provider())); err != nil {
	// 	logz.Error("[Watcher.applyChanges]", "err", err)
	// }
	if ev.Has(OpCreate) {
		logz.Debug("debug create")
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			s.setKV(key, val, true)
			logz.Debug("created: ", key, s.MustGet(key), "event", ev.Op())
		}
	} else if ev.Has(OpWrite) {
		logz.Debug("debug write")
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			s.setKV(key, val, false)
			logz.Debug("modified: ", key, s.MustGet(key), "event", ev.Op())
		}
	} else if ev.Has(OpRename) {
		logz.Debug("debug rename")
		for {
			key, _, ok := ev.Next()
			if !ok {
				break
			}
			// s.Set(key, val)
			logz.Debug("renamed: ", key, s.MustGet(key), "event", ev.Op())
		}
	} else if ev.Has(OpRemove) {
		logz.Debug("debug remove")
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			s.Remove(key)
			logz.Debug("removed: ", key, val, "event", ev.Op())
		}
	} else if ev.Has(OpChmod) {
		logz.Debug("debug chmod")
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			// s.Set(key, nil)
			logz.Debug("chmod: ", key, val, "event", ev.Op())
		}
	}
}
