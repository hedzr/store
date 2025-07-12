package hcl

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

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

// Unmarshal parses the given hcl bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	if err = hclsimple.Decode(
		"example.hcl", b,
		nil, &data,
	); err != nil {
		return
	}

	// err = hcl.Unmarshal(b, &data)

	// var o *ast.File

	// o, err = hcl.Parse(string(b))
	// if err != nil {
	// 	return
	// }

	// if err = hcl.DecodeObject(&data, o); err != nil {
	// 	return
	// }

	if p.flattenSlices {
		flattenHCL(data)
	}
	return
}

func ctyValue(v any) (ret cty.Value, err error) {
	switch z := v.(type) {
	// case map[string]any:
	// 	bazBlock := body.AppendNewBlock(k, nil)
	// 	bazBody := bazBlock.Body()
	// 	fill(bazBody, z)
	case fmt.Stringer:
		ret = cty.StringVal(z.String())
	case bool:
		ret = cty.BoolVal(z)
	case string:
		ret = cty.StringVal(z)
	case int:
		ret = cty.NumberIntVal(int64(z))
	case int8:
		ret = cty.NumberIntVal(int64(z))
	case int16:
		ret = cty.NumberIntVal(int64(z))
	case int32:
		ret = cty.NumberIntVal(int64(z))
	case int64:
		ret = cty.NumberIntVal(int64(z))
	case uint:
		ret = cty.NumberUIntVal(uint64(z))
	case uint8:
		ret = cty.NumberUIntVal(uint64(z))
	case uint16:
		ret = cty.NumberUIntVal(uint64(z))
	case uint32:
		ret = cty.NumberUIntVal(uint64(z))
	case uint64:
		ret = cty.NumberUIntVal(uint64(z))
	default:
		rv := reflect.ValueOf(z)
		if rv.Kind() == reflect.Slice {
			if l := rv.Len(); l > 0 {
				var values []cty.Value
				for ix := 0; ix < l; ix++ {
					iv := rv.Index(ix)
					if val, err := ctyValue(iv.Interface()); err == nil {
						values = append(values, val)
					} else {
						err = fmt.Errorf("[hcl-write][ctyValue] unprocessed v: v = %+v\n", v)
					}
				}
				ret = cty.ListVal(values)
			} else {
				ret = cty.ListValEmpty(cty.Type{})
			}
		} else {
			err = fmt.Errorf("[hcl-write][ctyValue] unprocessed v: v = %+v\n", v)
		}
	}
	return
}

func fill(p *ldr, body *hclwrite.Body, k string, v any) {
	switch z := v.(type) {
	case map[string]any:
		bazBlock := body.AppendNewBlock(k, nil)
		bazBody := bazBlock.Body()
		p.fill(bazBody, z)

	default:
		val, err := ctyValue(z)
		if err != nil {
			err = fmt.Errorf("[hcl-write][fill] unprocessed k: k = %v, err: %w\n", k, err)
		}
		body.SetAttributeValue(k, val)
	}
}

func (p *ldr) fill(body *hclwrite.Body, m map[string]any) {
	for k, v := range m {
		fill(p, body, k, v)
	}
}

// Marshal marshals the given config map to hcl bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	p.fill(rootBody, m)

	// rootBody.SetAttributeValue("string", cty.StringVal("bar")) // this is overwritten later
	// rootBody.AppendNewline()
	// rootBody.SetAttributeValue("object", cty.ObjectVal(map[string]cty.Value{
	// 	"foo": cty.StringVal("foo"),
	// 	"bar": cty.NumberIntVal(5),
	// 	"baz": cty.True,
	// }))
	// rootBody.SetAttributeValue("string", cty.StringVal("foo"))
	// rootBody.SetAttributeValue("bool", cty.False)
	// rootBody.SetAttributeTraversal("path", hcl.Traversal{
	// 	hcl.TraverseRoot{
	// 		Name: "env",
	// 	},
	// 	hcl.TraverseAttr{
	// 		Name: "PATH",
	// 	},
	// })
	// rootBody.AppendNewline()
	// fooBlock := rootBody.AppendNewBlock("foo", nil)
	// fooBody := fooBlock.Body()
	// rootBody.AppendNewBlock("empty", nil)
	// rootBody.AppendNewline()
	// barBlock := rootBody.AppendNewBlock("bar", []string{"a", "b"})
	// barBody := barBlock.Body()

	// fooBody.SetAttributeValue("hello", cty.StringVal("world"))

	// bazBlock := barBody.AppendNewBlock("baz", nil)
	// bazBody := bazBlock.Body()
	// bazBody.SetAttributeValue("foo", cty.NumberIntVal(10))
	// bazBody.SetAttributeValue("beep", cty.StringVal("boop"))
	// bazBody.SetAttributeValue("baz", cty.ListValEmpty(cty.String))

	data = f.Bytes()
	err = nil
	// err = store.ErrNotImplemented
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
