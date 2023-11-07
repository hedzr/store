package yaml

import (
	"os"

	yamlv3 "gopkg.in/yaml.v3"

	"github.com/hedzr/store"
)

func New() store.Codec {
	return &ldr{}
}

type ldr struct{}

// Unmarshal parses the given YAML bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = yamlv3.Unmarshal(b, &data)
	return
}

// Marshal marshals the given config map to YAML bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	return yamlv3.Marshal(m)
}

func (l *ldr) Load(file string) (data map[string]any, err error) {
	var f *os.File

	// data = make(map[string]any)

	if f, err = os.Open(file); err != nil {
		return
	}
	dec := yamlv3.NewDecoder(f)
	if err = dec.Decode(data); err != nil {
		return
	}
	return
}
