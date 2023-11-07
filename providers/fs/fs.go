//go:build go1.16
// +build go1.16

package fs

import (
	"io"
	"io/fs"

	"github.com/hedzr/store"
)

func New(fs fs.FS, pathname string, opts ...Opt) *pvdr {
	s := &pvdr{FS: fs, path: pathname}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

type Opt func(s *pvdr)
type pvdr struct {
	fs.FS
	path   string
	prefix string
	codec  store.Codec
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

func (s *pvdr) GetCodec() (codec store.Codec) { return s.codec }
func (s *pvdr) GetPosition() (pos string)     { return s.prefix }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.prefix = prefix }

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

func (s *pvdr) Value(key string) (value interface{}, ok bool) {
	ok = false
	return
}

func (s *pvdr) MustValue(key string) (value interface{}) {
	return
}

func (s *pvdr) Reader() (r *store.Reader, err error) {
	err = store.NotImplemented
	return
}

// Read returns the loaded map[string]interface{}.
func (s *pvdr) Read() (data map[string]interface{}, err error) {
	err = store.NotImplemented
	return
}

// ReadBytes is not supported by the confmap provider.
func (s *pvdr) ReadBytes() (data []byte, err error) {
	var f fs.File
	f, err = s.Open(s.path)
	if err != nil {
		return
	}
	defer f.Close()

	return io.ReadAll(f)
}

func (s *pvdr) Write(data []byte) (err error) {
	err = store.NotImplemented
	return
}
