package store

import (
	"context"
	"encoding/json"
	"io"
	"strconv"
	"sync/atomic"
	"time"

	logz "github.com/hedzr/logg/slog"
	"gopkg.in/hedzr/errors.v3"

	"github.com/hedzr/store/internal/radix"
)

func (s *storeS) inLoading() bool { return atomic.LoadInt32(&s.loading) == 1 }

func (s *storeS) WithinLoading(fn func()) {
	if atomic.CompareAndSwapInt32(&s.loading, 0, 1) {
		defer func() { atomic.CompareAndSwapInt32(&s.loading, 1, 0) }()
		fn()
	}
}

func (s *storeS) Load(ctx context.Context, opts ...LoadOpt) (wr Writeable, err error) {
	if atomic.CompareAndSwapInt32(&s.loading, 0, 1) {
		defer func() { atomic.CompareAndSwapInt32(&s.loading, 1, 0) }()

		loader := newLoader(s, opts...)

		var data map[string]any
		data, err = loader.tryLoad(ctx) // load dataset from source via loader
		if err != nil {
			return
		}

		// merge dataset into store
		if err = loader.loadMap(data, loader.Prefix(), true); err != nil {
			return
		}

		wr = loader

		loader.startWatch(ctx, loader)
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

//

//

//

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

func (s *loadS) tryLoad(ctx context.Context) (data map[string]any, err error) {
	if s.provider == nil {
		return
	}

	var b []byte

	// try Read() at first
	data, err = s.provider.Read()

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
				s.setKV(k, fp.MustValue(k), true)
			}
		}
	}
	if err != nil {
		return
	}

	// Decode it after loaded
	if s.codec != nil {
		data, err = s.codec.Unmarshal(b)
	} else if data == nil {
		// or fallback to json decoder
		err = json.Unmarshal(b, &data)
	}
	return
}

func (s *loadS) Save(ctx context.Context) (err error) { return s.trySave(ctx) }
func (s *loadS) trySave(ctx context.Context) (err error) {
	if s.codec != nil && s.provider != nil {
		var m map[string]any
		if m, err = s.GetM("", WithFilter[any](func(node radix.Node[any]) bool {
			return node.Modified() // && !strings.HasPrefix(node.Key(), "app.cmd.")
		})); err == nil {
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
