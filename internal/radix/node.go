package radix

import (
	"fmt"
	"strings"

	evendeep "github.com/hedzr/evendeep"

	"github.com/hedzr/is/term/color"
)

type nodeType int

const (
	NTBranch nodeType = iota // non-leaf nodes in a tree
	NTLeaf
	NTData = 1 << iota
	NTMask = NTLeaf
)

type nodeS[T any] struct {
	path        []rune // path fragment for this node.
	pathS       string // full path for performance
	children    []*nodeS[T]
	data        T
	description string
	comment     string
	tag         any
	nType       nodeType
}

type Extractor func(outputPtr any, defaultValue ...any) (err error)

func (s *nodeS[T]) isBranch() bool      { return s.nType&NTMask == NTBranch }
func (s *nodeS[T]) hasData() bool       { return s.nType&NTData != 0 }
func (s *nodeS[T]) Description() string { return s.description }
func (s *nodeS[T]) Comment() string     { return s.comment }
func (s *nodeS[T]) Tag() any            { return s.tag }

func (s *nodeS[T]) Data() (data T) {
	if !s.isBranch() {
		data = s.data
	}
	return
}

func (s *nodeS[T]) endsWith(ch rune) bool {
	if len(s.path) == 0 {
		return false
	}
	return s.path[len(s.path)-1] == ch
}

func (s *nodeS[T]) endsWithLite(ch rune) bool {
	return s.path[len(s.path)-1] == ch
}

func (s *nodeS[T]) remove(item *nodeS[T]) (removed bool) {
	if item == nil {
		return
	}
	for i, c := range s.children {
		if c == item {
			removed, s.children = true, append(s.children[:i], s.children[i+1:]...)
			break
		}
	}
	return
}

func (s *nodeS[T]) findCommonPrefixLength(word []rune) (length int) {
	ml := min(len(word), len(s.path))
	for length < ml && word[length] == s.path[length] {
		length++
	}
	return
}

func (s *nodeS[T]) insertInternal(word []rune, fullPath string, data T) (oldData any) {
	if strings.Contains(string(word), " ") {
		word = []rune(strings.ReplaceAll(string(word), " ", "-"))
		fullPath = strings.ReplaceAll(fullPath, " ", "-")
	}

	base, ourLen, wordLen := s, len(s.path), len(word)
	if ourLen == 0 {
		if wordLen > 0 && len(s.children) == 0 {
			_ = base.insertAsLeaf(word, fullPath, data)
			return
		}
	}

	// var newNode *nodeS[T]
	var cpl int
	if ourLen > 0 && wordLen > 0 {
		cpl = base.findCommonPrefixLength(word)
	}

	if cpl < ourLen {
		// eg: insert 'apple' into 'appZ', or insert 'appZ' into 'apple'
		base.split(cpl, word) // split this as 'app' and 'Z'/'le'
	}

	if cpl < wordLen {
		// eg: insert 'apple' into 'app'
		if cpl > 0 {
			word = word[cpl:]
		}
		matched, child := base.matchChildren(word)
		if matched {
			child.insert(word, fullPath, data)
		} else {
			base.insertAsLeaf(word, fullPath, data)
		}
	} else {
		// hit this node,
		base.nType |= NTData
		oldData, base.data = base.data, data
	}
	return
}

func (s *nodeS[T]) split(pos int, word []rune) (newNode *nodeS[T]) {
	// origPath, origPathS := s.path, s.pathS
	// if assertEnabled {
	// 	defer func() {
	// 		origPath, origPathS = []rune(origPathS), string(origPath)
	// 	}()
	// }

	tip("[store/radix] [split] original path, pathS: %q, %q", string(s.path), s.pathS)

	d := len(s.path) - pos

	newNode = &nodeS[T]{
		path:     s.path[pos:],
		pathS:    s.pathS, // [pos:], // s.pathS[len(s.pathS)-d:],
		children: s.children,
		data:     s.data,
		nType:    s.nType,
	}
	assert(strings.HasSuffix(newNode.pathS, string(newNode.path)), "newNode: pathS should end with path")

	s.path = s.path[:pos]
	s.pathS = s.pathS[:len(s.pathS)-d] // s.pathS[:pos] //
	s.children = []*nodeS[T]{newNode}
	s.nType = NTBranch
	assert(strings.HasSuffix(s.pathS, string(s.path)), "parentNode: pathS(%q) should end with path(%q)", s.pathS, string(s.path))
	return
}

func (s *nodeS[T]) insertAsLeaf(word []rune, fullPath string, data T) (newNode *nodeS[T]) {
	newNode = &nodeS[T]{
		path:  word,
		pathS: fullPath,
		nType: NTLeaf | NTData,
		data:  data,
	}
	assert(strings.HasSuffix(newNode.pathS, string(newNode.path)), "newNode: pathS should end with path")
	s.children = append(s.children, newNode)
	return
}

func (s *nodeS[T]) matchChildren(word []rune) (matched bool, child *nodeS[T]) {
	for _, child = range s.children {
		if child.path[0] == word[0] {
			matched = true
			break
		}
	}
	return
}

func (s *nodeS[T]) matchR(word []rune, delimiter rune, parentNode *nodeS[T]) (matched, partialMatched bool, child, parent *nodeS[T]) {
	wl, l := len(word), len(s.path)
	if wl == 0 {
		return true, false, s, parentNode
	}

	matchedL, minL := 0, min(l, wl)
	for ; matchedL < minL; matchedL++ {
		if s.path[matchedL] == word[matchedL] {
			continue
		}
		if matchedL < l {
			// partial matched.
			if matchedL < wl {
				// eg: matching 'apple' on 'apk'
				return false, false, nil, nil
			}
			// eg: matching 'app' on 'apple', or 'apk' on 'apple'
			return true, false, s, parentNode
		}
	}

	if matchedL == l-1 && s.path[matchedL] == delimiter {
		matched, child, parent = true, s, parentNode
	} else if minL < l && matchedL == minL {
		partialMatched, child, parent = true, s, parentNode
	} else if minL >= l && matchedL == minL {
		matched, child, parent = true, s, parentNode
	}
	if minL < wl {
		for _, child = range s.children {
			matched, partialMatched, child, parent = child.matchR(word[minL:], delimiter, s)
			if matched || partialMatched {
				return
			}
		}
	}
	return
}

func (s *nodeS[T]) dump(noColor bool) string {
	var sb strings.Builder
	return s.dumpR(&sb, 0, noColor)
}

const col1Width = 32
const branchTitle = "<B>"
const leafTitle = "<L>"

func (s *nodeS[T]) dumpR(sb *strings.Builder, lvl int, noColor bool) string {
	sb.WriteString(strings.Repeat("  ", lvl))
	if len(s.path) == 0 {
		if lvl > 0 {
			sb.WriteString("(nil)\n")
		}
	} else {
		sb.WriteString(string(s.path))
		if col1Width-lvl*2-len(s.path) > 0 {
			sb.WriteString(strings.Repeat(" ", col1Width-lvl*2-len(s.path)))
		} else {
			sb.WriteByte(' ')
		}

		if noColor {
			if s.isBranch() {
				sb.WriteString(branchTitle)
			} else {
				sb.WriteString(leafTitle)
			}
		} else {
			if s.isBranch() {
				sb.WriteString(color.ToDim(branchTitle))
			} else {
				sb.WriteString(color.ToDim(leafTitle))
			}
		}

		if s.hasData() {
			sb.WriteString(" ")
			sb.WriteString(s.pathS)
			sb.WriteString(" => ")
			sb.WriteString(color.ToDim(fmt.Sprint(s.data)))
		}

		if !strings.HasSuffix(s.pathS, string(s.path)) {
			sb.WriteString(fmt.Sprintf(" [WRONG path & pathS: %q / %q]", string(s.path), s.pathS))
		}
		sb.WriteByte('\n')
	}

	for _, child := range s.children {
		child.dumpR(sb, lvl+1, noColor)
	}
	return sb.String()
}

func (s *nodeS[T]) Dup() (newNode *nodeS[T]) {
	newNode = &nodeS[T]{
		path:  s.path,
		pathS: s.pathS,
		nType: s.nType,
	}

	newNode.children = make([]*nodeS[T], 0, len(s.children))
	for _, ch := range s.children {
		newNode.children = append(newNode.children, ch.Dup())
	}

	data := evendeep.MakeClone(s.data)
	switch z := data.(type) {
	case T:
		newNode.data = z
	case *T:
		newNode.data = *z
	}
	return
}

func (s *nodeS[T]) Walk(cb func(prefix, key string, node *nodeS[T])) {
	s.walk(0, cb)
}

func (s *nodeS[T]) walk(level int, cb func(prefix, key string, node *nodeS[T])) {
	cb(s.pathS, string(s.path), s)
	for _, ch := range s.children {
		ch.walk(level+1, cb)
	}
}
