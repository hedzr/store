package radix

import (
	"strconv"

	"gopkg.in/hedzr/errors.v3"
)

func (s *trieS[T]) Prefix() string { return s.prefix }

func (s *trieS[T]) WithPrefix(prefix string) (entry Trie[T]) {
	return s.withPrefixR(prefix)
}

func (s *trieS[T]) withPrefixR(prefix string) (entry *trieS[T]) {
	return &trieS[T]{root: s.root, prefix: s.join(s.prefix, prefix), delimiter: s.delimiter}
}

func (s *trieS[T]) WithPrefixReplaced(prefix string) (entry Trie[T]) {
	return s.withPrefixSimple(prefix)
}

func (s *trieS[T]) withPrefixSimple(prefix string) (entry *trieS[T]) {
	return &trieS[T]{root: s.root, prefix: prefix, delimiter: s.delimiter}
}

func (s *trieS[T]) SetPrefix(prefix string) {
	s.prefix = prefix
}

func (s *trieS[T]) loadMap(m map[string]any) (err error) {
	ec := errors.New()
	defer ec.Defer(&err)
	for k, v := range m {
		s.loadMapByValueType(ec, m, k, v)
	}
	return
}

func (s *trieS[T]) loadMapByValueType(ec errors.Error, m map[string]any, k string, v any) {
	switch vv := v.(type) {
	case map[string]any:
		ec.Attach(s.withPrefixSimple(k).loadMap(vv))
	case []map[string]any:
		buf := make([]byte, 0, len(k)+16)
		for i, mm := range vv {
			buf = append(buf, k...)
			buf = append(buf, byte(s.delimiter))
			buf = strconv.AppendInt(buf, int64(i), 10)
			ec.Attach(s.withPrefixR(string(buf)).loadMap(mm))
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
