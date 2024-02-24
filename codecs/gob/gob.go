package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/hedzr/store"
)

// New makes a new instance for Gob encoder and decoder.
//
// For serializing a struct with unexported fields, you need
// to implement [gob.GobEncoder] and [gob.GobDecoder]
// interfaces.
func New() store.Codec {
	return &ldr{}
}

type ldr struct {
	dec *gob.Decoder
	enc *gob.Encoder
}

// Unmarshal parses the given bytes.
func (p *ldr) Unmarshal(b []byte) (data map[string]any, err error) {
	if p.dec == nil {
		r := bytes.NewReader(b)
		p.dec = gob.NewDecoder(r)
	}
	err = p.dec.Decode(data)
	return
}

// Marshal marshals the given config map to bytes.
func (p *ldr) Marshal(m map[string]any) (data []byte, err error) {
	if p.enc == nil {
		sb := bytes.NewBuffer(data)
		p.enc = gob.NewEncoder(sb)
	}
	err = p.enc.Encode(m)
	return
}
