package consul

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/hedzr/logg/slog"

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

	// onUpdated OnUpdated
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

// func WithOnUpdated(cb OnUpdated) Opt {
// 	return func(s *pvdr) {
// 		s.onUpdated = cb
// 	}
// }
//
// type OnUpdated func()

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

func (s *pvdr) Read() (data map[string]store.ValPkg, err error) {
	var kv = s.Client.KV()
	var pairs api.KVPairs
	var pair *api.KVPair

	data = make(map[string]store.ValPkg)

	if s.recursive {
		pairs, _, err = kv.List(s.position, nil)
		if err != nil {
			return
		}

		// Detailed information can be obtained via hedzr/store.GetString(key), which key are:
		//
		// "parent1.CreateIndex"
		// "parent1.Flags"
		// "parent1.LockIndex"
		// "parent1.ModifyIndex"
		// "parent1.Session"
		// "parent1.Value"
		if s.processMeta {
			for _, pair = range pairs {
				m := make(map[string]any)
				m["CreateIndex"] = strconv.FormatUint(pair.CreateIndex, 10)
				m["Flags"] = strconv.FormatUint(pair.Flags, 10)
				m["LockIndex"] = strconv.FormatUint(pair.LockIndex, 10)
				m["ModifyIndex"] = strconv.FormatUint(pair.ModifyIndex, 10)

				if pair.Session == "" {
					m["Session"] = "-"
				} else {
					m["Session"] = fmt.Sprintf("%s", pair.Session)
				}

				m["Value"] = string(pair.Value)

				data[s.NormalizeKey(pair.Key)] = store.ValPkg{
					Value:   string(pair.Value),
					Desc:    "",
					Comment: "",
					Tag:     m,
				}
			}
		} else {
			for _, pair := range pairs {
				data[s.NormalizeKey(pair.Key)] = store.ValPkg{
					Value: string(pair.Value),
				}
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
		m["CreateIndex"] = strconv.FormatUint(pair.CreateIndex, 10)
		m["Flags"] = strconv.FormatUint(pair.Flags, 10)
		m["LockIndex"] = strconv.FormatUint(pair.LockIndex, 10)
		m["ModifyIndex"] = strconv.FormatUint(pair.ModifyIndex, 10)

		if pair.Session == "" {
			m["Session"] = "-"
		} else {
			m["Session"] = fmt.Sprintf("%s", pair.Session)
		}

		m["Value"] = string(pair.Value)

		data[s.NormalizeKey(pair.Key)] = store.ValPkg{
			Value:   string(pair.Value),
			Desc:    "",
			Comment: "",
			Tag:     m,
		}
	} else {
		data[s.NormalizeKey(pair.Key)] = store.ValPkg{
			Value: string(pair.Value),
		}
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
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) Write(data []byte) (err error) {
	err = store.ErrNotImplemented // todo implement consul writer
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

// Watch watches for changes in the Consul API and triggers a callback.
func (s *pvdr) Watch(ctx context.Context, cb func(event any, err error)) (err error) {
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
		slog.Debug("consul watching plan has been started", "type", p["type"], "key", s.position, "addr", s.config.Address)

		// no need to try to shutdown from ctx triggering. because the
		// plan will be stopped while consul Client.Close() called.

		err = s.plan.Run(s.config.Address)
		if err != nil {
			slog.Error("consul watching plan has error", "err", err)
			return
		}
	}()

	return nil
}
