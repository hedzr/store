package file

import (
	"time"

	"github.com/hedzr/store"
)

type changeS struct {
	realPath string
	idx      int

	lastOp        store.Op
	lastEvent     string
	lastEventTime time.Time

	provider store.Provider
}

func (s *changeS) Path() string             { return s.realPath }
func (s *changeS) Op() store.Op             { return s.lastOp }
func (s *changeS) Has(op store.Op) bool     { return uint64(s.lastOp)&uint64(op) != 0 }
func (s *changeS) Timestamp() time.Time     { return s.lastEventTime }
func (s *changeS) Provider() store.Provider { return s.provider }
func (s *changeS) Next() (key string, val any, ok bool) {
	if s.idx == 0 {
		key, val, ok = s.realPath, s.realPath, true
		s.idx++
	}
	return
}
func (s *changeS) Set() {
	s.idx = 0
}
