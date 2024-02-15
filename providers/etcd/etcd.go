package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/consul/api/watch"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/hedzr/store"
)

func New(opts ...Opt) *pvdr {
	s := &pvdr{config: clientv3.Config{DialTimeout: 5 * time.Second}}
	for _, opt := range opts {
		opt(s)
	}
	_ = s.prepare()
	return s
}

type Opt func(s *pvdr)
type pvdr struct {
	*clientv3.Client
	config clientv3.Config
	// endpoints   []string
	// dialTimeout time.Duration
	watchEnabled               bool
	codec                      store.Codec
	storePrefix                string
	prefixOrKey                bool
	stripPrefix, prependPrefix string
	delimiter                  string
	limit                      bool
	maxLimit                   int64
	processMeta                bool
}

func WithCodec(codec store.Codec) Opt {
	return func(s *pvdr) {
		s.codec = codec
	}
}

func WithPosition(prefix string) Opt {
	return func(s *pvdr) {
		s.storePrefix = prefix
	}
}

func WithWatchEnabled(b bool) Opt {
	return func(s *pvdr) {
		s.watchEnabled = b
	}
}

func WithProcessMeta(b bool) Opt {
	return func(s *pvdr) {
		s.processMeta = b
	}
}

func WithRecursive(b bool) Opt {
	return func(s *pvdr) {
		s.prefixOrKey = b
	}
}

func WithEndpoints(peers ...string) Opt {
	return func(s *pvdr) {
		s.config.Endpoints = peers
	}
}

func WithDialTimeout(timeout time.Duration) Opt {
	return func(s *pvdr) {
		s.config.DialTimeout = timeout
	}
}

func WithEtcdConfig(config clientv3.Config) Opt {
	return func(s *pvdr) {
		s.config = config
	}
}

func WithStorePrefix(prefix string) Opt {
	return func(s *pvdr) {
		s.storePrefix, s.prefixOrKey = prefix, true
	}
}

func WithDelimiter(d string) Opt {
	return func(s *pvdr) {
		s.delimiter = d
	}
}

func WithStripPrefix(prefix string) Opt {
	return func(s *pvdr) {
		s.stripPrefix = prefix
	}
}

func WithPrependPrefix(prefix string) Opt {
	return func(s *pvdr) {
		s.prependPrefix = prefix
	}
}

func WithKey(key string) Opt {
	return func(s *pvdr) {
		s.storePrefix, s.prefixOrKey = key, false
	}
}

//

func (s *pvdr) prepare() (err error) {
	s.Client, err = clientv3.New(s.config)
	return
}

func (s *pvdr) Count() int {
	return 0
}

func (s *pvdr) Has(key string) bool {
	return false
}

func (s *pvdr) Next() (key string, eol bool) {
	eol = true
	return
}

func (s *pvdr) Keys() (keys []string, err error) {
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) Value(key string) (value any, ok bool) {
	ok = false
	return
}

func (s *pvdr) MustValue(key string) (value any) {
	return
}

func (s *pvdr) Reader() (r store.Reader, err error) {
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) NormalizeKey(key string) string {
	if s.stripPrefix != "" {
		key = strings.TrimPrefix(key, s.stripPrefix+"/")
	}
	if s.prependPrefix != "" {
		key = strings.Join([]string{s.prependPrefix, key}, s.delimiter)
	}
	if s.delimiter != "" {
		key = strings.ReplaceAll(key, "/", s.delimiter)
	}
	return key
}

func (s *pvdr) Read() (data map[string]any, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.DialTimeout)
	defer cancel()

	var resp *clientv3.GetResponse
	if s.prefixOrKey {
		if s.limit {
			resp, err = s.Get(ctx, s.storePrefix, clientv3.WithPrefix(), clientv3.WithLimit(s.maxLimit))
			if err != nil {
				return
			}
		} else {
			resp, err = s.Get(ctx, s.storePrefix, clientv3.WithPrefix())
			if err != nil {
				return
			}
		}
	} else {
		resp, err = s.Get(ctx, s.storePrefix)
		if err != nil {
			return
		}
	}

	data = make(map[string]any, len(resp.Kvs))
	for _, r := range resp.Kvs {
		data[s.NormalizeKey(string(r.Key))] = string(r.Value)
	}

	return
}

func (s *pvdr) ReadBytes() (data []byte, err error) {
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) Write(data []byte) (err error) {
	err = store.ErrNotImplemented // todo implement etcd writer
	return
}

func (s *pvdr) GetCodec() (codec store.Codec) { return s.codec }
func (s *pvdr) GetPosition() (pos string)     { return s.storePrefix }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.storePrefix = prefix }

func (s *pvdr) Close() {
	if s.Client != nil {
		s.Client.Close()
	}
}

type changeS struct {
	ev  *clientv3.Event
	idx int

	lastOp        store.Op
	lastEventTime time.Time

	plan *watch.Plan

	provider store.Provider
}

func (s *changeS) Path() string             { return "" }
func (s *changeS) Op() store.Op             { return s.lastOp }
func (s *changeS) Has(op store.Op) bool     { return uint64(s.lastOp)&uint64(op) != 0 }
func (s *changeS) Timestamp() time.Time     { return s.lastEventTime }
func (s *changeS) Provider() store.Provider { return s.provider }
func (s *changeS) Next() (key string, val any, ok bool) {
	if s.idx == 0 {
		key, val, ok = s.provider.(*pvdr).NormalizeKey(string(s.ev.Kv.Key)), s.ev.Kv.Value, true
		s.idx++
	}
	return
}
func (s *changeS) Set(ev *clientv3.Event) {
	s.lastEventTime = time.Now()
	s.ev = ev
	s.lastOp = store.OpNone
	if ev.IsCreate() {
		s.lastOp |= store.OpCreate
	}
	if ev.IsModify() {
		s.lastOp |= store.OpWrite
	}
}

// Watch watches for changes in the Consul API and triggers a callback.
func (s *pvdr) Watch(cb func(event any, err error)) error {
	if s.watchEnabled == false {
		return nil
	}

	var w clientv3.WatchChan
	var lastChange = changeS{provider: s}

	go func() {
		if s.prefixOrKey {
			w = s.Client.Watch(context.Background(), s.storePrefix, clientv3.WithPrefix())
		} else {
			w = s.Client.Watch(context.Background(), s.storePrefix)
		}

		for wresp := range w {
			for _, ev := range wresp.Events {
				lastChange.Set(ev)
				cb(&lastChange, nil)
			}
		}
	}()

	return nil
}
