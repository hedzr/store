package radix

import (
	"strconv"

	"gopkg.in/hedzr/errors.v3"
)

// Prefix returns the current prefix setting.
func (s *trieS[T]) Prefix() string { return s.prefix }

// WithPrefix makes a new Trie instance with a new prefix,
// which is the joint value with the current prefix setting
// and the given prefix value.
func (s *trieS[T]) WithPrefix(prefix ...string) (entry Trie[T]) {
	return s.withPrefixImpl(prefix...)
}

func (s *trieS[T]) withPrefixImpl(prefix ...string) (entry *trieS[T]) {
	return s.dupS(s.root, s.join1(s.prefix, prefix...))
}

// WithPrefixReplaced makes a new Trie instance with a new
// prefix with the given prefix value.
// The current prefix setting was ignored.
func (s *trieS[T]) WithPrefixReplaced(newPrefix ...string) (entry Trie[T]) {
	return s.withPrefixReplacedImpl(newPrefix...)
}

func (s *trieS[T]) withPrefixReplacedImpl(newPrefix ...string) (entry *trieS[T]) {
	return s.dupS(s.root, s.join(newPrefix...))
}

// SetPrefix replaces the current prefix setting with the given new value.
func (s *trieS[T]) SetPrefix(newPrefix ...string) {
	s.prefix = s.join(newPrefix...)
}

func (s *trieS[T]) loadMap(m map[string]any) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	for k, v := range m {
		s.loadMapByValueType(ec, m, k, v)
	}
	return
}

func (s *trieS[T]) loadMapByValueType(ec errors.Error, m map[string]any, k string, v any) { //nolint:revive,unparam
	switch vv := v.(type) {
	case map[string]any:
		ec.Attach(s.withPrefixReplacedImpl(k).loadMap(vv))
	case []map[string]any:
		buf := make([]byte, 0, len(k)+16)
		for i, mm := range vv {
			buf = append(buf, k...)
			buf = append(buf, byte(s.delimiter))
			buf = strconv.AppendInt(buf, int64(i), 10)
			ec.Attach(s.withPrefixImpl(string(buf)).loadMap(mm))
			buf = buf[:0]
		}
	case []any:
		buf := make([]byte, 0, len(k)+16)
		for i, mm := range vv {
			if s.prefix != "" {
				buf = append(buf, s.prefix...)
				buf = append(buf, byte(s.delimiter))
			}
			buf = append(buf, k...)
			buf = append(buf, byte(s.delimiter))
			buf = strconv.AppendInt(buf, int64(i), 10)
			s.loadMapByValueType(ec, m, string(buf), mm)
			buf = buf[:0]
		}
	default:
		s.Set(k, v.(T))
	}
}
