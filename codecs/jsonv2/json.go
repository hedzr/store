package json

import (
	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"os"

	"github.com/hedzr/store"
)

func New(opts ...Opt) store.Codec {
	s := &ldr{pretty: true}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithPretty(pretty bool) Opt {
	return func(s *ldr) {
		s.pretty = pretty
	}
}

func WithFlattenSlice(flattenSlice bool) Opt {
	return func(s *ldr) {
		s.flattenSlice = flattenSlice
	}
}

type (
	Opt func(s *ldr)
	ldr struct {
		flattenSlice bool
		pretty       bool
	}
)

var _ store.Codec = (*ldr)(nil)

// Unmarshal parses the given JSON bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = json.Unmarshal(b, &data,
		jsonv1.DefaultOptionsV1(),
		jsontext.AllowDuplicateNames(false),
	)
	return
}

// Marshal marshals the given config map to JSON bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	if p.pretty {
		return jsonv1.MarshalIndent(m, "", "  ")
	}
	return json.Marshal(m)
}

func (l *ldr) Load(file string) (data map[string]any, err error) {
	var f *os.File
	if f, err = os.Open(file); err != nil {
		return
	}

	dec := jsontext.NewDecoder(f)
	if err = json.UnmarshalDecode(dec, &data,
		jsonv1.DefaultOptionsV1(),
		jsontext.AllowDuplicateNames(false),
	); err != nil {
		return
	}

	return
}
