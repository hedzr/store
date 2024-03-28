//go:build delve
// +build delve

package radix

func (s *nodeS[T]) insert(word []rune, fullPath string, data T) (node Node[T], oldData any) {
	node, oldData = s.insertInternal(word, fullPath, data)
	str := s.dump(true) // check integrity
	_ = str
	return
}
