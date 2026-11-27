package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/hedzr/store"
)

type changeS struct {
	ev  *clientv3.Event
	idx int

	lastOp        store.Op
	lastEventTime time.Time

	// plan *watch.Plan

	provider store.Provider
}

func (s *changeS) Path() string             { return "" }
func (s *changeS) Op() store.Op             { return s.lastOp }
func (s *changeS) Has(op store.Op) bool     { return uint64(s.lastOp)&uint64(op) != 0 }
func (s *changeS) Timestamp() time.Time     { return s.lastEventTime }
func (s *changeS) Provider() store.Provider { return s.provider }
func (s *changeS) Next() (key string, val any, ok bool) {
	if s.idx == 0 {
		key, val, ok = s.provider.(*pvdr).NormalizeKey(string(s.ev.Kv.Key)), s.ev.Kv.Value, true
		if b, ok1 := val.([]byte); ok1 {
			val = string(b)
		}
		s.idx++
	}
	return
}
func (s *changeS) Set(ev *clientv3.Event) {
	s.lastEventTime = time.Now()
	s.ev = ev
	s.lastOp = store.OpNone
	s.idx = 0
	if ev.IsCreate() {
		s.lastOp |= store.OpCreate
	}
	if ev.IsModify() {
		s.lastOp |= store.OpWrite
	}
	if ev.Type == clientv3.EventTypeDelete && ev.Kv.CreateRevision != ev.Kv.ModRevision {
		s.lastOp |= store.OpRemove
	}
}
