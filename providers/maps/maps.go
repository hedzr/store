package maps

import (
	"github.com/hedzr/store"
	"github.com/hedzr/store/cvt"
)

func New(m map[string]any, delimiter string, opts ...Opt) *pvdr { //nolint:revive
	s := &pvdr{}
	for _, opt := range opts {
		opt(s)
	}

	cp := cvt.Normalize(cvt.Copy(m), nil)
	if s.delimiter != "" {
		cp = cvt.Deflate(cp, s.delimiter)
	}
	s.m = cp

	return s
}

type Opt func(s *pvdr)

type pvdr struct {
	m         map[string]any
	delimiter string
	prefix    string
	codec     store.Codec
}

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
func (s *pvdr) Read() (data map[string]any, err error) {
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
