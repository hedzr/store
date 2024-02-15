package hcl

import (
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"

	"github.com/hedzr/store"
)

func New(opts ...Opt) store.Codec {
	s := &ldr{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithFlattenSlices(b bool) Opt {
	return func(s *ldr) {
		s.flattenSlices = b
	}
}

type Opt func(s *ldr)
type ldr struct{ flattenSlices bool }

// Unmarshal parses the given YAML bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = hcl.Unmarshal(b, &data)

	var o *ast.File

	o, err = hcl.Parse(string(b))
	if err != nil {
		return
	}

	if err = hcl.DecodeObject(&data, o); err != nil {
		return
	}

	if p.flattenSlices {
		flattenHCL(data)
	}
	return
}

// Marshal marshals the given config map to YAML bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	// return hcl.Marshal(m)
	err = store.ErrNotImplemented
	return
}

func (l *ldr) Load(file string) (data map[string]any, err error) {
	// var f *os.File

	// // data = make(map[string]any)

	// if f, err = os.Open(file); err != nil {
	// 	return
	// }
	// dec := yamlv3.NewDecoder(f)
	// if err = dec.Decode(data); err != nil {
	// 	return
	// }

	err = store.ErrNotImplemented
	return
}

// flattenHCL flattens an unmarshalled HCL structure where maps
// turn into slices -- https://github.com/hashicorp/hcl/issues/162.
func flattenHCL(mp map[string]any) {
	for k, val := range mp {
		if v, ok := val.([]map[string]any); ok {
			if len(v) == 1 {
				mp[k] = v[0]
			}
		}
	}
	for _, val := range mp {
		if v, ok := val.(map[string]any); ok {
			flattenHCL(v)
		}
	}
}
