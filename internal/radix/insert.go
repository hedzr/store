//go:build !delve
// +build !delve

package radix

func (s *nodeS[T]) insert(word []rune, fullPath string, data T) (node Node[T], oldData any) {
	return s.insertInternal(word, fullPath, data)
}
