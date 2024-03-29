package maps

import (
	"context"
	"sync/atomic"

	"github.com/hedzr/store"
	"github.com/hedzr/store/internal/cvt"
)

// New makes a new instance for importing a map.
//
// The map must be a map[string]any, it can be nested.
func New(m map[string]any, delimiter string, opts ...Opt) *pvdr { //nolint:revive
	s := &pvdr{delimiter: delimiter}
	for _, opt := range opts {
		opt(s)
	}

	cp := cvt.Normalize(cvt.Copy(m), nil)
	if s.delimiter != "" {
		cp = cvt.Deflate(cp, s.delimiter)
	}

	// s.m = cp

	s.m = make(map[string]store.ValPkg)
	for k, v := range cp {
		s.m[k] = store.ValPkg{
			Value:   v,
			Desc:    "",
			Comment: "",
			Tag:     nil,
		}
	}
	return s
}

type Opt func(s *pvdr)

func WithCodec(codec store.Codec) Opt {
	return func(s *pvdr) {
		s.codec = codec
	}
}

func WithPosition(prefix string) Opt {
	return func(s *pvdr) {
		s.prefix = prefix
	}
}

func WithDelimiter(d string) Opt {
	return func(s *pvdr) {
		s.delimiter = d
	}
}

type pvdr struct {
	m         map[string]store.ValPkg
	delimiter string
	prefix    string
	codec     store.Codec
	watching  int32
}

func (s *pvdr) GetCodec() (codec store.Codec) { return s.codec }
func (s *pvdr) GetPosition() (pos string)     { return s.prefix }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.prefix = prefix }

func (s *pvdr) Count() int {
	return 0
}

func (s *pvdr) Has(key string) bool { //nolint:revive
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

func (s *pvdr) Value(key string) (value any, ok bool) { //nolint:revive
	ok = false
	return
}

func (s *pvdr) MustValue(key string) (value any) { //nolint:revive
	return
}

func (s *pvdr) Reader() (r store.Reader, err error) { //nolint:revive
	err = store.ErrNotImplemented
	return
}

// Read returns the loaded map[string]interface{}.
func (s *pvdr) Read() (data map[string]store.ValPkg, err error) {
	return s.m, nil
}

// ReadBytes is not supported by the confmap provider.
func (s *pvdr) ReadBytes() (data []byte, err error) {
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) Write(data []byte) (err error) { //nolint:revive
	err = store.ErrNotImplemented
	return
}

func (s *pvdr) Close() {
	atomic.CompareAndSwapInt32(&s.watching, 1, 0)
}

func (s *pvdr) Watch(ctx context.Context, cb func(event any, err error)) (err error) {
	if !atomic.CompareAndSwapInt32(&s.watching, 0, 1) {
		return
	}

	// todo do some stuff here to enable watching for maps provider

	return
}
