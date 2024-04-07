package radix

import (
	"fmt"
	"strings"

	"github.com/hedzr/evendeep"
)

type nodeType int

const (
	NTBranch   nodeType    = iota // non-leaf nodes in a tree
	NTLeaf                        // leaf node
	NTData     = 1 << iota        // node has data field, only if it is a leaf node
	NTModified                    // node attrs(data, desc, comment, or tag) modified?
	NTMask     = NTLeaf           // mask for checking if it's a branch or leaf
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

var _ Node[any] = (*nodeS[any])(nil) // assertion helper

type Extractor func(outputPtr any, defaultValue ...any) (err error) // data field extractor

func (s *nodeS[T]) isBranch() bool      { return s.nType&NTMask == NTBranch }
func (s *nodeS[T]) hasData() bool       { return s.nType&NTData != 0 }
func (s *nodeS[T]) Modified() bool      { return s.nType&NTModified != 0 } // modification state
func (s *nodeS[T]) Description() string { return s.description }           // description field
func (s *nodeS[T]) Comment() string     { return s.comment }               // comment field
func (s *nodeS[T]) Tag() any            { return s.tag }                   // tag field
func (s *nodeS[T]) Key() string         { return s.pathS }                 // key field is the full path of this node

func (s *nodeS[T]) IsLeaf() bool  { return s.nType&NTMask == NTLeaf } // leaf node?
func (s *nodeS[T]) HasData() bool { return s.nType&NTData != 0 }      //nolint:revive //data field is valid?

// SetModified sets the modified state to true or false.
//
// To clear the state, using ResetModified;
// Or flip the state with ToggleModified.
func (s *nodeS[T]) SetModified(b bool) { //nolint:revive
	if b {
		s.nType |= NTModified
	} else {
		s.nType &= ^NTModified
	}
}

// ToggleModified flips the modified state.
func (s *nodeS[T]) ToggleModified() {
	s.nType ^= NTModified
	// if s.Modified() {
	// 	s.nType &= ^NTModified
	// } else {
	// 	s.nType |= NTModified
	// }
}

// ResetModified clears the modified state.
func (s *nodeS[T]) ResetModified() {
	s.nType &= ^NTModified
}

// Data returns the Data field of a node.
func (s *nodeS[T]) Data() (data T) {
	if !s.isBranch() {
		data = s.data
	}
	return
}

// SetData sets the Data field of a node.
func (s *nodeS[T]) SetData(data T) {
	if !s.isBranch() {
		s.data = data
	}
}

// SetComment sets the Description and Comment field.
func (s *nodeS[T]) SetComment(desc, comment string) { //nolint:revive
	if s.isBranch() {
		return
	}
	s.description, s.comment = desc, comment
}

// SetTag sets the Tag field.
//
// You may save any value into a Tag field.
func (s *nodeS[T]) SetTag(tag any) { //nolint:revive
	if s.isBranch() {
		return
	}
	s.tag = tag
}

func (s *nodeS[T]) endsWith(ch rune) bool { //nolint:revive
	if len(s.path) == 0 {
		return false
	}
	return s.path[len(s.path)-1] == ch
}

func (s *nodeS[T]) endsWithLite(ch rune) bool { //nolint:revive,unused
	return s.path[len(s.path)-1] == ch
}

func (s *nodeS[T]) remove(item *nodeS[T]) (removed bool) { //nolint:revive
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

func (s *nodeS[T]) insertInternal(word []rune, fullPath string, data T) (node *nodeS[T], oldData any) { //nolint:revive
	if strings.Contains(string(word), " ") {
		word = []rune(strings.ReplaceAll(string(word), " ", "-")) //nolint:revive
		fullPath = strings.ReplaceAll(fullPath, " ", "-")         //nolint:revive
	}

	base, ourLen, wordLen := s, len(s.path), len(word)
	if ourLen == 0 {
		if wordLen > 0 && len(s.children) == 0 {
			node = base.insertAsLeaf(word, fullPath, data)
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
			word = word[cpl:] //nolint:revive
		}
		matched, child := base.matchChildren(word)
		if matched {
			var n Node[T]
			n, oldData = child.insert(word, fullPath, data)
			if nn, ok := n.(*nodeS[T]); ok {
				node = nn
			}
		} else {
			node = base.insertAsLeaf(word, fullPath, data)
		}
	} else {
		// hit this node,
		base.nType |= NTData
		node, oldData, base.data = base, base.data, data
	}
	return
}

func (s *nodeS[T]) split(pos int, word []rune) (newNode *nodeS[T]) { //nolint:unparam,revive
	// origPath, origPathS := s.path, s.pathS
	// if assertEnabled {
	// 	defer func() {
	// 		origPath, origPathS = []rune(origPathS), string(origPath)
	// 	}()
	// }

	tip("[store/radix] [split] original path, pathS: %q, %q", string(s.path), s.pathS)

	d := len(s.path) - pos
	_ = word

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

func (s *nodeS[T]) matchR(word []rune, delimiter rune, parentNode *nodeS[T]) (matched, partialMatched bool, child, parent *nodeS[T]) { //nolint:revive
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
	} else if minL >= l && matchedL == minL && minL > 0 && matchedL > 0 {
		matched, child, parent = true, s, parentNode
	}
	if minL < wl {
		if len(s.children) == 0 {
			matched, partialMatched = false, true
			return
		}
		for _, child = range s.children {
			matched, partialMatched, child, parent = child.matchR(word[minL:], delimiter, s)
			if matched || partialMatched {
				return
			}
		}
	}
	return
}

func (s *nodeS[T]) dump(noColor bool) string { //nolint:revive
	var sb strings.Builder
	return s.dumpR(&sb, 0, noColor)
}

const (
	col1Width   = 32
	branchTitle = "<B>"
	leafTitle   = "<L>"
)

func iif[T any](cond bool, tv, fv T) (rv T) {
	if cond {
		return tv
	}
	return fv
}

func (s *nodeS[T]) dumpR(sb *strings.Builder, lvl int, noColor bool) string { //nolint:revive
	_, _ = sb.WriteString(strings.Repeat("  ", lvl))
	if len(s.path) == 0 {
		if lvl > 0 {
			_, _ = sb.WriteString("(nil)\n")
		}
	} else {
		_, _ = sb.WriteString(string(s.path))
		if width := col1Width - lvl*2 - len(s.path); width > 0 {
			_, _ = sb.WriteString(strings.Repeat(" ", width))
		} else {
			_ = sb.WriteByte(' ')
		}

		isbr := s.isBranch()

		if noColor {
			_, _ = sb.WriteString(iif(isbr, branchTitle, leafTitle))
			// if s.isBranch() {
			// 	_, _ = sb.WriteString(branchTitle)
			// } else {
			// 	_, _ = sb.WriteString(leafTitle)
			// }
		} else {
			_, _ = sb.WriteString(ColorToDim(iif(isbr, branchTitle, leafTitle)))
			// if s.isBranch() {
			// 	_, _ = sb.WriteString(ColorToDim(branchTitle))
			// } else {
			// 	_, _ = sb.WriteString(ColorToDim(leafTitle))
			// }
		}

		if s.hasData() {
			_, _ = sb.WriteString(" ")
			_, _ = sb.WriteString(s.pathS)
			_, _ = sb.WriteString(" => ")
			_, _ = sb.WriteString(ColorToDim(fmt.Sprint(s.data)))
		}

		if s.comment != "" {
			_, _ = sb.WriteString(ColorToColor(FgLightGreen, " // "+s.comment))
		}

		if s.tag != nil {
			_, _ = sb.WriteString(" | tag = ")
			_, _ = sb.WriteString(ColorToColor(FgGreen, fmt.Sprint(s.tag)))
		}

		if s.description != "" {
			_, _ = sb.WriteString(ColorToColor(FgLightGreen, " ~ "+s.description))
		}

		if !strings.HasSuffix(s.pathS, string(s.path)) {
			_, _ = fmt.Fprintf(sb, " [WRONG path & pathS: %q / %q]", string(s.path), s.pathS)
		}
		_ = sb.WriteByte('\n')
	}

	for _, child := range s.children {
		child.dumpR(sb, lvl+1, noColor)
	}
	return sb.String()
}

// Dup or Clone makes an exact deep copy of this node and all of its children.
//
// The cost of cloneing a large tree or node is expensive.
func (s *nodeS[T]) Dup() (newNode *nodeS[T]) { //nolint:revive
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
	case *T:
		newNode.data = *z
	case T:
		newNode.data = z
	}
	return
}

// Walk navigates all of chilren from this node.
func (s *nodeS[T]) Walk(cb func(path, fragment string, node Node[T])) { //nolint:revive
	s.walk(0, cb)
}

func (s *nodeS[T]) walk(level int, cb func(path, fragment string, node Node[T])) { //nolint:revive
	cb(s.pathS, string(s.path), s)
	for _, ch := range s.children {
		ch.walk(level+1, cb)
	}
}
