package env

import (
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hedzr/store"
)

func New(opts ...Opt) *pvdr {
	s := &pvdr{lowerCase: true, underline2dot: false}
	for _, opt := range opts {
		opt(s)
	}
	_ = s.prepare()
	return s
}

type pvdr struct {
	codec         store.Codec
	prefix        string
	stripped      string
	storePrefix   string
	cb            func(key string) string
	keys          []string
	m             map[string]any
	pos           int
	lowerCase     bool
	underline2dot bool
}

type Opt func(s *pvdr)

func WithCodec(codec store.Codec) Opt {
	return func(s *pvdr) {
		s.codec = codec
	}
}

// WithPrefix gives a filter prefix for env var name.
//
// Only names which have given prefix are available.
// Such as app.Name()+"_".
//
// The filtered names will be stripped by strippedPrefix.
func WithPrefix(prefix string, strippedPrefix ...string) Opt {
	return func(s *pvdr) {
		s.prefix = strings.ToUpper(prefix)
		for _, ss := range strippedPrefix {
			s.stripped = ss
		}
	}
}

// WithStorePrefix gives a dotted key prefix for store.Store.
//
// Such as: "app.cmd", ..
func WithStorePrefix(position string) Opt {
	return func(s *pvdr) {
		s.storePrefix = position
	}
}

func WithKeyCB(cb func(key string) string) Opt {
	return func(s *pvdr) {
		s.cb = cb
	}
}

func WithLowerCase(b ...bool) Opt {
	return func(s *pvdr) {
		var lc = true
		for _, bb := range b {
			lc = bb
		}
		s.lowerCase = lc
	}
}

func WithUnderlineToDot(b ...bool) Opt {
	return func(s *pvdr) {
		var lc = true
		for _, bb := range b {
			lc = bb
		}
		s.underline2dot = lc
	}
}

func (s *pvdr) prepare() (err error) {
	re := regexp.MustCompile(`([^_]*)_([^_])`)
	vec := os.Environ()
	s.m = make(map[string]any, len(vec))
	s.keys = make([]string, 0, len(vec))
	for _, p := range vec {
		pos := strings.Index(p, "=")
		if pos > 0 {
			k, v := p[:pos], p[pos+1:]
			if s.prefix != "" {
				if !strings.HasPrefix(k, s.prefix) {
					continue
				}
			}
			if s.lowerCase {
				k = strings.ToLower(k)
			}
			if s.stripped != "" {
				k = strings.TrimPrefix(k, s.stripped)
			}
			if s.underline2dot {
				k = k[:1] + re.ReplaceAllString(k[1:], "$1.$2")
			}
			if s.cb != nil {
				k = s.cb(k)
			}
			if s.storePrefix != "" {
				k = s.storePrefix + "." + k
			}
			s.m[k] = v
			s.keys = append(s.keys, k)
		}
	}
	sort.Strings(s.keys)
	s.pos = 0
	return
}

func (s *pvdr) Count() int {
	return len(s.keys)
}

func (s *pvdr) Has(key string) bool {
	_, ok := s.m[key]
	return ok
}

func (s *pvdr) Next() (key string, eol bool) {
	if eol = s.pos < len(s.keys); !eol {
		key = s.keys[s.pos]
		s.pos++
	}
	return
}

func (s *pvdr) Keys() (keys []string, err error) {
	keys = s.keys
	return
}

func (s *pvdr) Value(key string) (value any, ok bool) {
	value, ok = s.m[key]
	return
}

func (s *pvdr) MustValue(key string) (value any) {
	value = s.m[key]
	return
}

func (s *pvdr) Reader() (r *store.Reader, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) Read() (data map[string]any, err error) {
	data = s.m
	return
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
func (s *pvdr) GetPosition() (pos string)     { return s.prefix }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.storePrefix = prefix }
