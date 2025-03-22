package store

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"sync/atomic"
	"time"

	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/evendeep"
	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store/radix"
)

func (s *storeS) inLoading() bool { return atomic.LoadInt32(&s.loading) == 1 }

// WithinLoading is a helper to 'load' a 'fn'. The 'fn' will be
// run as is, and the internal flag 's.loading' will be set at
// beginning of fn executing, and reset at ending of fn.
func (s *storeS) WithinLoading(fn func()) {
	if atomic.CompareAndSwapInt32(&s.loading, 0, 1) {
		defer func() { atomic.CompareAndSwapInt32(&s.loading, 1, 0) }()
		fn()
	}
}

// Load loads an external data source by the specified Provider,
// a Codec parser is optional.
//
// WithProvider and WithCodec are useful. The sample code is:
//
//	s := newBasicStore()
//	if _, err := s.Load(
//	   context.TODO(),
//	   store.WithStorePrefix("app.json"),
//	   store.WithCodec(json.New()),
//	   store.WithProvider(file.New("../testdata/4.json")),
//
//	   store.WithStoreFlattenSlice(true),
//	); err != nil {
//	   t.Fatalf("failed: %v", err)
//	}
func (s *storeS) Load(ctx context.Context, opts ...LoadOpt) (wr Writeable, err error) { //nolint:revive
	if atomic.CompareAndSwapInt32(&s.loading, 0, 1) {
		defer func() { atomic.CompareAndSwapInt32(&s.loading, 1, 0) }()

		loader := newLoader(s, opts...)

		var data map[string]ValPkg
		var bin map[string]any
		data, bin, err = loader.tryLoad(ctx) // load dataset from source via loader
		if err != nil {
			return
		}

		// merge dataset into store
		prefix := loader.Prefix()
		ok := false
		if data != nil {
			if err = loader.loadMapDedicated(data, prefix, true); err != nil {
				return
			}
			ok = true
		}
		if bin != nil {
			if err = loader.loadMap(bin, prefix, true, nil); err != nil {
				return
			}
			ok = true
		}

		if ok {
			wr = loader
			loader.startWatch(ctx, loader)
		}
	}
	return
}

// func (s *storeS) Save(ctx context.Context, wr Writeable, opts ...LoadOpt) (err error) {
// 	if atomic.CompareAndSwapInt32(&s.saving, 0, 1) {
// 		defer func() { atomic.CompareAndSwapInt32(&s.saving, 1, 0) }()
//
// 		loader := wr // newLoader(s, opts...)
//
// 		err = loader.Save(ctx)
// 	}
// 	return
// }

type lmOnSet func(node radix.Node[any])

func (s *storeS) loadMapDedicated(m map[string]ValPkg, position string, creating bool) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	for k, v := range m {
		s.loadMapByValueType(ec, position, k, v.Value, creating, func(node radix.Node[any]) {
			node.SetComment(v.Desc, v.Comment)
			node.SetTag(v.Tag)
		})
	}
	return
}

func (s *storeS) loadMapAny(m map[any]any, position string, creating bool, onSet lmOnSet) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	cvt := evendeep.Cvt{}
	for k, v := range m {
		s.loadMapByValueType(ec, position, cvt.String(k), v, creating, onSet)
	}
	return
}

func (s *storeS) loadMap(m map[string]any, position string, creating bool, onSet lmOnSet) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	for k, v := range m {
		s.loadMapByValueType(ec, position, k, v, creating, onSet)
	}
	return
}

func privateSetter(ss *storeS, position, k string, v any, creating bool, onSet lmOnSet) {
	set := ss.WithPrefixReplaced(position).(*storeS)
	defer func() { atomic.StoreInt32(&set.loading, 0) }()
	set.setKV(k, v, creating, onSet)
}

func (s *storeS) loadMapByValueType(ec errors.Error, position, k string, v any, creating bool, onSet lmOnSet) { //nolint:revive
	switch vv := v.(type) {
	case ValPkg:
		s.loadMapByValueType(ec, position, k, vv.Value, creating, onSet)
	case map[string]any:
		ec.Attach(s.loadMap(vv, s.join(position, k), creating, onSet))
	case map[any]any:
		ec.Attach(s.loadMapAny(vv, s.join(position, k), creating, onSet))
	case []map[string]any:
		if s.flattenSlice {
			buf := make([]byte, 0, len(k)+16)
			for i, mm := range vv {
				buf = append(buf, k...)
				buf = append(buf, byte(s.Delimiter()))
				buf = strconv.AppendInt(buf, int64(i), 10)
				ec.Attach(s.loadMap(mm, s.join(position, string(buf)), creating, onSet))
				buf = buf[:0]
			}
			break
		}

		privateSetter(s, position, k, v, creating, onSet)

		// if cc, ok := s.WithPrefixReplaced(position).(interface {
		// 	setKV(path string, data any, createOrModify bool, onSet lmOnSet) (node radix.Node[any], oldData any)
		// }); ok {
		// 	cc.setKV(k, v, creating, onSet)
		// }
	case []map[any]any:
		if s.flattenSlice {
			buf := make([]byte, 0, len(k)+16)
			for i, mm := range vv {
				buf = append(buf, k...)
				buf = append(buf, byte(s.Delimiter()))
				buf = strconv.AppendInt(buf, int64(i), 10)
				ec.Attach(s.loadMapAny(mm, s.join(position, string(buf)), creating, onSet))
				buf = buf[:0]
			}
			break
		}

		privateSetter(s, position, k, v, creating, onSet)

		// if cc, ok := s.WithPrefixReplaced(position).(interface {
		// 	setKV(path string, data any, createOrModify bool, onSet lmOnSet) (node radix.Node[any], oldData any)
		// }); ok {
		// 	cc.setKV(k, v, creating, onSet)
		// }
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
				s.loadMapByValueType(ec, position, string(buf), mm, creating, onSet)
				buf = buf[:0]
			}
			break
		}

		privateSetter(s, position, k, v, creating, onSet)

		// if cc, ok := s.WithPrefixReplaced(position).(interface {
		// 	setKV(path string, data any, createOrModify bool, onSet lmOnSet) (node radix.Node[any], oldData any)
		// }); ok {
		// 	cc.setKV(k, v, creating, onSet)
		// }
	default:
		privateSetter(s, position, k, v, creating, onSet)

		// if cc, ok := set.(interface {
		// 	setKV(path string, data any, createOrModify bool, onSet lmOnSet) (node radix.Node[any], oldData any)
		// }); ok {
		// 	cc.setKV(k, v, creating, onSet)
		// }
	}
}

// Watchable tips that a Provider can watch its external data source
type Watchable interface {
	// Watch accepts user's func and callback it when the external
	// data source is changing, creating or deleting.
	//
	// The supported oprations are specified in Op.
	//
	// Tne user's func checks 'event' for which operation was occurring.
	// For more info, see also storeS.Load, storeS.applyExternalChanges,
	// and loader.startWatch.
	Watch(ctx context.Context, cb func(event any, err error)) error

	// Close provides a closer to cleanup the peripheral gracefully
	Close()
	// basics.Peripheral
}

// Change is an abstract interface for Watchable object.
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

func (s *storeS) startWatch(ctx context.Context, loader *loadS) {
	if !s.allowWatch || loader.provider == nil {
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

func (s *storeS) applyChanges(ev Change) { //nolint:revive
	// if err := s.Load(WithProvider(ev.Provider())); err != nil {
	// 	logz.Error("[Watcher.applyChanges]", "err", err)
	// }
	if hasCreate, hasWrite := ev.Has(OpCreate), ev.Has(OpWrite); hasCreate || hasWrite {
		logz.Debug("debug create/write", "create", hasCreate, "write", hasWrite)
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			s.setKV(key, val, hasCreate, nil)
			logz.Debug("created/wrote: ", key, s.MustGet(key), "event", ev.Op())
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
	} else if hasRename, hasChmod := ev.Has(OpRename), ev.Has(OpChmod); hasRename || hasChmod {
		logz.Debug("debug rename/chmod", "rename", hasRename, "chmod", hasChmod)
		for {
			key, val, ok := ev.Next()
			if !ok {
				break
			}
			// s.Set(key, nil)
			logz.Debug("renamed/chmod'ed: ", key, val, "event", ev.Op())
		}
	}
}

//

//

//

func newLoader(st *storeS, opts ...LoadOpt) *loadS {
	loader := &loadS{
		storeS:   st,
		codec:    nil,
		provider: nil,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(loader)
		}
	}

	if loader.provider != nil {
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
	}

	// loader.provider.WithStorePrefix(s.prefix)
	return loader
}

type loadS struct {
	*storeS
	position string
	codec    Codec
	provider Provider
}

type LoadOpt func(*loadS) // options for loadS

// WithProvider is commonly required. It specify what Provider
// will be [storeS.Load].
func WithProvider(provider Provider) LoadOpt {
	return func(s *loadS) {
		s.provider = provider
	}
}

// WithCodec specify the decoder to decode the loaded data.
func WithCodec(codec Codec) LoadOpt {
	return func(s *loadS) {
		s.codec = codec
	}
}

// WithStorePrefix gives a prefix position, which is the store
// location that the external settings will be merged at.
func WithStorePrefix(prefix string) LoadOpt {
	return func(s *loadS) {
		s.storeS = s.storeS.WithPrefixReplaced(prefix).(*storeS)
	}
}

// WithPosition sets the
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
// For example, if a store has a prefix 'app.server',
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
// By using the default setting, i.e., keepPrefix == false, we will strip
// the may-be-there prefix if necessary. So both "app.server.tls" and "tls"
// will work properly as you really want.
func WithKeepPrefix[T any](b bool) radix.MOpt[T] {
	return radix.WithKeepPrefix[T](b)
}

// WithFilter can be used in calling GetM(path, ...)
func WithFilter[T any](filter radix.FilterFn[T]) radix.MOpt[T] {
	return radix.WithFilter[T](filter)
}

// WithoutFlattenKeys allows returns a nested map.
// If the keys contain delimiter char, they will be split as
// nested sub-map.
func WithoutFlattenKeys[T any](b bool) radix.MOpt[T] {
	return radix.WithoutFlattenKeys[T](b)
}

// tryLoad inspect the provider's api, try reading settings in the best way.
//
// See also [storeS.Load].
func (s *loadS) tryLoad(ctx context.Context) (data map[string]ValPkg, bin map[string]any, err error) { //nolint:revive
	if s.provider == nil {
		return
	}

	_ = ctx

	// try Read() at first
	data, err = s.provider.Read()
	if err == nil {
		return // Read ok, return the data directly
	}

	var b []byte

	if errors.Is(err, ErrNotImplemented) {
		// the 2nd is OnceProvider and/or StreamProvider
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
				s.setKV(k, fp.MustValue(k), true, nil)
			}
		}
	}
	if err != nil {
		return
	}

	// Decode it after loaded
	if s.codec != nil {
		bin, err = s.codec.Unmarshal(b)
	} else if data == nil {
		// or fallback to json decoder
		err = json.Unmarshal(b, &bin)
	}

	return
}

func (s *loadS) Save(ctx context.Context) (err error) { return s.trySave(ctx) }
func (s *loadS) trySave(ctx context.Context) (err error) { //nolint:revive
	_ = ctx
	if s.codec != nil && s.provider != nil {
		// logz.InfoContext(ctx, "Write-Back", "position", s.position)
		var m map[string]any
		logz.DebugContext(ctx, "Write-Back checking", "src", s.provider)
		if m, err = s.GetM("",
			WithFilter[any](func(node radix.Node[any]) bool {
				return node.Modified() // && !strings.HasPrefix(node.Key(), "app.cmd.")
			}),
			// WithKeepPrefix[any](true),
			WithoutFlattenKeys[any](true),
		); err == nil && m != nil && len(m) > 0 {
			logz.DebugContext(ctx, "Write-Back checked and invoking", "src", s.provider)
			var data []byte
			if data, err = s.codec.Marshal(m); err == nil {
				switch fp := s.provider.(type) {
				case OnceProvider:
					err = fp.Write(data)
				default:
					err = ErrNotImplemented
				}

				if errors.Is(err, ErrNotImplemented) {
					if wr, ok := s.provider.(io.Writer); ok {
						_, err = wr.Write(data)
					}
				}
			}
		}
	}
	return
}
