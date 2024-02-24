package json

import (
	"encoding/json"
	"os"

	"github.com/hedzr/store"
)

func New() store.Codec {
	return &ldr{}
}

type ldr struct{}

// Unmarshal parses the given JSON bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	err = json.Unmarshal(b, &data)
	return
}

// Marshal marshals the given config map to JSON bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	return json.Marshal(m)
}

func (l *ldr) Load(file string) (data map[string]any, err error) {
	var f *os.File

	// data = make(map[string]any)

	if f, err = os.Open(file); err != nil {
		return
	}
	dec := json.NewDecoder(f)
	if err = dec.Decode(&data); err != nil {
		return
	}
	return
}
