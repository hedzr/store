package toml

import (
	"github.com/hedzr/store"

	"github.com/pelletier/go-toml"
)

func New() store.Codec {
	return &ldr{}
}

type ldr struct{}

// Unmarshal parses the given YAML bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = toml.Unmarshal(b, &data)
	return
}

// Marshal marshals the given config map to YAML bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	return toml.Marshal(m)
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

	err = store.NotImplemented
	return
}
