package consul

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/hedzr/store"
)

func New(opts ...Opt) *pvdr {
	s := &pvdr{recursive: true, processMeta: true}
	for _, opt := range opts {
		opt(s)
	}
	_ = s.prepare()
	return s
}

type Opt func(s *pvdr)
type pvdr struct {
	*api.Client
	config *api.Config

	watchEnabled          bool
	codec                 store.Codec
	position, storePrefix string
	recursive             bool
	processMeta           bool

	stripPrefix, prependPrefix string
	delimiter                  string // replace consul slash '/' with delimiter

	plan *watch.Plan
}

func WithCodec(codec store.Codec) Opt {
	return func(s *pvdr) {
		s.codec = codec
	}
}

func WithPosition(prefix string) Opt {
	return func(s *pvdr) {
		s.position = prefix
	}
}

func WithWatchEnabled(b bool) Opt {
	return func(s *pvdr) {
		s.watchEnabled = b
	}
}

func WithRecursive(b bool) Opt {
	return func(s *pvdr) {
		s.recursive = b
	}
}

func WithProcessMeta(b bool) Opt {
	return func(s *pvdr) {
		s.processMeta = b
	}
}

func WithPrefix(prefix string) Opt {
	return func(s *pvdr) {
		s.storePrefix = prefix
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

func WithDelimiter(delimiter string) Opt {
	return func(s *pvdr) {
		s.delimiter = delimiter
	}
}

func WithConsulConfig(cfg *api.Config) Opt {
	return func(s *pvdr) {
		s.config = cfg
	}
}

func (s *pvdr) prepare() (err error) {
	if s.config != nil {
		s.Client, err = api.NewClient(s.config)
	}
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
	err = store.NotImplemented
	return
}

func (s *pvdr) Value(key string) (value any, ok bool) {
	ok = false
	return
}

func (s *pvdr) MustValue(key string) (value any) {
	return
}

func (s *pvdr) Reader() (r *store.Reader, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) Read() (data map[string]any, err error) {
	var kv = s.Client.KV()
	var pairs api.KVPairs
	var pair *api.KVPair

	data = make(map[string]any)

	if s.recursive {
		pairs, _, err = kv.List(s.position, nil)
		if err != nil {
			return
		}

		// Detailed information can be obtained using standard koanf flattened delimited keys:
		// For example:
		// "parent1.CreateIndex"
		// "parent1.Flags"
		// "parent1.LockIndex"
		// "parent1.ModifyIndex"
		// "parent1.Session"
		// "parent1.Value"
		if s.processMeta {
			for _, pair := range pairs {
				m := make(map[string]any)
				m["CreateIndex"] = fmt.Sprintf("%d", pair.CreateIndex)
				m["Flags"] = fmt.Sprintf("%d", pair.Flags)
				m["LockIndex"] = fmt.Sprintf("%d", pair.LockIndex)
				m["ModifyIndex"] = fmt.Sprintf("%d", pair.ModifyIndex)

				if pair.Session == "" {
					m["Session"] = "-"
				} else {
					m["Session"] = fmt.Sprintf("%s", pair.Session)
				}

				m["Value"] = string(pair.Value)

				data[s.NormalizeKey(pair.Key)] = m
			}
		} else {
			for _, pair := range pairs {
				data[s.NormalizeKey(pair.Key)] = string(pair.Value)
			}
		}

		return
	}

	pair, _, err = kv.Get(s.position, nil)
	if err != nil {
		return
	}

	if s.processMeta {
		m := make(map[string]any)
		m["CreateIndex"] = fmt.Sprintf("%d", pair.CreateIndex)
		m["Flags"] = fmt.Sprintf("%d", pair.Flags)
		m["LockIndex"] = fmt.Sprintf("%d", pair.LockIndex)
		m["ModifyIndex"] = fmt.Sprintf("%d", pair.ModifyIndex)

		if pair.Session == "" {
			m["Session"] = "-"
		} else {
			m["Session"] = fmt.Sprintf("%s", pair.Session)
		}

		m["Value"] = string(pair.Value)

		data[s.NormalizeKey(pair.Key)] = m
	} else {
		data[s.NormalizeKey(pair.Key)] = string(pair.Value)
	}

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

func (s *pvdr) ReadBytes() (data []byte, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) Write(data []byte) (err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) GetCodec() (codec store.Codec) { return s.codec }
func (s *pvdr) GetPosition() (pos string)     { return s.position }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.position = prefix }

func (s *pvdr) Close() {
	if s.plan != nil && !s.plan.IsStopped() {
		s.plan.Stop()
	}
}

type changeS struct {
	val   any
	idx   int
	index uint64 // consul action index

	lastOp        store.Op
	lastEvent     string
	lastEventTime time.Time

	plan *watch.Plan

	provider store.Provider
}

// func (s *changeS) Key() string              { return string(s.ev.Kv.Key) }
// func (s *changeS) Val() any                 { return s.val }

func (s *changeS) Path() string             { return s.val.(string) }
func (s *changeS) Op() store.Op             { return s.lastOp }
func (s *changeS) Has(op store.Op) bool     { return uint64(s.lastOp)&uint64(op) != 0 }
func (s *changeS) Timestamp() time.Time     { return s.lastEventTime }
func (s *changeS) Provider() store.Provider { return s.provider }
func (s *changeS) Next() (key string, val any, ok bool) {
	if s.provider.(*pvdr).recursive {
		if pairs, yes := s.val.(api.KVPairs); yes {
		RetryNextPair:
			if s.idx < len(pairs) {
				kvp := pairs[s.idx]
				if s.lastOp == store.OpRemove || s.checkOpIndex(kvp) {
					key, val, ok = s.provider.(*pvdr).NormalizeKey(kvp.Key), kvp.Value, true
				} else {
					s.idx++
					goto RetryNextPair
				}
				s.idx++
			}
		}
	} else {
		if kvp, yes := s.val.(*api.KVPair); yes {
			if s.idx == 0 {
				if s.checkOpIndex(kvp) {
					key, val, ok = s.provider.(*pvdr).NormalizeKey(kvp.Key), kvp.Value, true
				}
				s.idx++
			}
		}
	}
	return
}

func (s *changeS) Set(val any, idx uint64) (df func()) {
	s.lastEventTime = time.Now()
	s.index, s.idx = idx, 0
	df = s.findOpIndex(idx, val)
	if df == nil {
		df = func() {}
	}
	return
}

func (s *changeS) findOpIndex(idx uint64, val any) (df func()) {
	if s.provider.(*pvdr).recursive {
		if pairs, ok := val.(api.KVPairs); ok {
			for ix := 0; ix < len(pairs); ix++ {
				kvp := pairs[ix]
				if s.checkOpIndex(kvp) {
					s.val = val
					return
				}
			}

			s.lastOp = store.OpRemove
			df = func() { s.val = val }
			var removed api.KVPairs
			if s.val != nil {
				for _, p := range s.val.(api.KVPairs) {
					var found bool
					for _, q := range pairs {
						if p.Key == q.Key {
							found = true
							break
						}
					}
					if !found {
						removed = append(removed, p)
					}
				}
			}
			s.val = removed
		}
	} else {
		if kvp, ok := val.(*api.KVPair); ok {
			if !s.checkOpIndex(kvp) {
				s.lastOp = store.OpRemove
				df = func() { s.val = val }
				return
			}
			s.val = val
			return
		}
	}
	return
}

func (s *changeS) checkOpIndex(pair *api.KVPair) bool {
	if pair.CreateIndex == s.index {
		s.lastOp = store.OpCreate
		return true
	} else if pair.ModifyIndex == s.index {
		s.lastOp = store.OpWrite
		return true
	}
	return false
}

// Watch watches for changes in the Consul API and triggers a callback.
func (s *pvdr) Watch(cb func(event any, err error)) (err error) {
	if s.watchEnabled == false {
		return nil
	}

	p := make(map[string]any)

	if s.recursive {
		p["type"] = "keyprefix"
		p["prefix"] = s.position
	} else {
		p["type"] = "key"
		p["key"] = s.position
	}

	var lastChange = changeS{provider: s}

	lastChange.plan, err = watch.Parse(p)
	if err != nil {
		return err
	}

	s.plan = lastChange.plan
	s.plan.Handler = func(idx uint64, val any) {
		defer lastChange.Set(val, idx)()
		cb(&lastChange, nil)
	}
	// s.plan.HybridHandler = func(bpv watch.BlockingParamVal, val interface{}) {
	// 	lastChange.Set(val)
	// 	cb(lastChange, nil)
	// }

	go func() {
		s.plan.Run(s.config.Address)
	}()

	return nil
}
