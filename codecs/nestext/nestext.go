package nestext

import (
	"bytes"
	"errors"

	"github.com/npillmayer/nestext"
	"github.com/npillmayer/nestext/ntenc"

	"github.com/hedzr/store"
)

func New() store.Codec {
	return &ldr{}
}

type ldr struct{}

// Unmarshal parses the given NestedText bytes.
//
// If the NT content does not reflect a dict (NT allows
// for top-level lists or strings as well), the content
// will be wrapped into a dict with a single key
// named "nestedtext".
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	// err = nestext.Unmarshal(b, &data)

	var ok bool
	var result interface{}
	result, err = nestext.Parse(bytes.NewReader(b), nestext.TopLevel("dict"))
	if err != nil {
		return
	}

	data, ok = result.(map[string]interface{})
	if !ok {
		err = errors.New("NestedText configuration expected to be a dict at top-level")
	}
	return
}

// Marshal marshals the given config map to NestedText bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	// return nestext.Marshal(m)

	var buf bytes.Buffer
	_, err = ntenc.Encode(m, &buf)
	if err == nil {
		data = buf.Bytes()
	}
	return
}

func (l *ldr) Load(file string) (data map[string]any, err error) {
	// var f *os.File
	//
	// // data = make(map[string]any)
	//
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
