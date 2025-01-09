//go:build delve
// +build delve

package radix

func (s *nodeS[T]) insert(word []rune, fullPath string, data T, trie *trieS[T], cb OnSetEx[T]) (node Node[T], oldData any) {
	node, oldData = s.insertInternal(word, fullPath, data, trie, cb)
	str := s.dump(true) // check integrity
	_ = str
	return
}
