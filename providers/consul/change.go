package consul

import (
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/hedzr/store"
)

type changeS struct {
	val   any
	idx   int
	index uint64 // consul action index

	lastOp        store.Op
	lastEvent     string
	lastEventTime time.Time

	plan     *watch.Plan
	provider store.Provider
}

// func (s *changeS) Key() string              { return string(s.ev.Kv.Key) }
// func (s *changeS) Val() any                 { return s.val }

func (s *changeS) Path() string             { return s.val.(string) }
func (s *changeS) Op() store.Op             { return s.lastOp }
func (s *changeS) Has(op store.Op) bool     { return uint64(s.lastOp)&uint64(op) != 0 }
func (s *changeS) Timestamp() time.Time     { return s.lastEventTime }
func (s *changeS) Provider() store.Provider { return s.provider }
func (s *changeS) Next() (key string, val any, ok bool) {
	if s.provider.(*pvdr).recursive {
		if pairs, yes := s.val.(api.KVPairs); yes {
		RetryNextPair:
			if s.idx < len(pairs) {
				kvp := pairs[s.idx]
				if s.lastOp == store.OpRemove || s.checkOpIndex(kvp) {
					key, val, ok = s.provider.(*pvdr).NormalizeKey(kvp.Key), kvp.Value, true
				} else {
					s.idx++
					goto RetryNextPair
				}
				s.idx++
			}
		}
	} else {
		if kvp, yes := s.val.(*api.KVPair); yes {
			if s.idx == 0 {
				if s.checkOpIndex(kvp) {
					key, val, ok = s.provider.(*pvdr).NormalizeKey(kvp.Key), kvp.Value, true
				}
				s.idx++
			}
		}
	}
	return
}

func (s *changeS) Set(val any, idx uint64) (df func()) {
	s.lastEventTime = time.Now()
	s.index, s.idx = idx, 0
	df = s.findOpIndex(idx, val)
	if df == nil {
		df = func() {}
	}
	return
}

func (s *changeS) findOpIndex(idx uint64, val any) (df func()) {
	if s.provider.(*pvdr).recursive {
		if pairs, ok := val.(api.KVPairs); ok {
			for ix := 0; ix < len(pairs); ix++ {
				kvp := pairs[ix]
				if s.checkOpIndex(kvp) {
					s.val = val
					return
				}
			}

			s.lastOp = store.OpRemove
			df = func() { s.val = val }
			var removed api.KVPairs
			if s.val != nil {
				for _, p := range s.val.(api.KVPairs) {
					var found bool
					for _, q := range pairs {
						if p.Key == q.Key {
							found = true
							break
						}
					}
					if !found {
						removed = append(removed, p)
					}
				}
			}
			s.val = removed
		}
	} else {
		if kvp, ok := val.(*api.KVPair); ok {
			if !s.checkOpIndex(kvp) {
				s.lastOp = store.OpRemove
				df = func() { s.val = val }
				return
			}
			s.val = val
			return
		}
	}
	return
}

func (s *changeS) checkOpIndex(pair *api.KVPair) bool {
	if pair.CreateIndex == s.index {
		s.lastOp = store.OpCreate
		return true
	} else if pair.ModifyIndex == s.index {
		s.lastOp = store.OpWrite
		return true
	}
	return false
}
