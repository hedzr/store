package hjson

import (
	"github.com/hjson/hjson-go/v4"

	"github.com/hedzr/store"
)

func New() store.Codec {
	return &ldr{flattenSlice: true}
}

type ldr struct {
	flattenSlice bool
}

var _ store.Codec = (*ldr)(nil)

// Unmarshal parses the given YAML bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = hjson.Unmarshal(b, &data)
	return
}

// Marshal marshals the given config map to YAML bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	return hjson.Marshal(m)
}

func (p *ldr) MarshalEx(m map[string]store.ValPkg) (data []byte, err error) {
	return
}

func (p *ldr) UnmarshalEx(b []byte) (data map[string]store.ValPkg, err error) {
	var node hjson.Node
	if err = hjson.UnmarshalWithOptions(b, &node, hjson.DefaultDecoderOptions()); err != nil {
		return
	}

	data = make(map[string]store.ValPkg)
	err = p.decode(data, &node)
	return
}

func (p *ldr) decode(data map[string]store.ValPkg, node *hjson.Node) (err error) {
	if node.Value == nil {
		return
	}
	if s, ok := node.Value.(string); ok { // a bad case here
		cm := node.Cm
		data[s] = store.ValPkg{Value: s,
			Desc: cm.Before, Comment: cm.Key + cm.After, Tag: nil}
		return
	}

	for i := 0; i < node.Len(); i++ {
		k, v, e := node.AtIndex(i)
		if e != nil {
			return e
		}

		n := node.NI(i)
		cm := n.Cm

		switch v.(type) {
		case *hjson.OrderedMap:
			nm := make(map[string]store.ValPkg)
			err = p.decode(nm, n)
			data[k] = store.ValPkg{Value: nm,
				Desc: cm.Before, Comment: cm.Key + cm.After, Tag: nil}
			continue
		case []any:
			if !p.flattenSlice {
				data[k] = store.ValPkg{Value: v,
					Desc: cm.Before, Comment: cm.Key + cm.After, Tag: nil}
				continue
			}

			nm := make(map[string]store.ValPkg)
			// k = strconv.Itoa(i)
			data[k] = store.ValPkg{Value: nm,
				Desc: cm.Before, Comment: cm.Key + cm.After, Tag: nil}
			err = p.decode(nm, n)
			continue
		default:
			data[k] = store.ValPkg{Value: v,
				Desc: cm.Before, Comment: cm.Key + cm.After, Tag: nil}
		}
	}
	return
}

func (p *ldr) Load(file string) (data map[string]any, err error) {
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
