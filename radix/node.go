package radix

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hedzr/evendeep"
	logz "github.com/hedzr/logg/slog"
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
	pathS       string // full path for performance, on a node we want to retrieve its full path as Key()
	children    []*nodeS[T]
	data        T
	description string
	comment     string
	tag         any
	nType       nodeType
	rw          *sync.RWMutex
}

var _ Node[any] = (*nodeS[any])(nil) // assertion helper

type Extractor func(outputPtr any, defaultValue ...any) (err error) // data field extractor

func (s *nodeS[T]) isBranch() bool      { return s.nType&NTMask == NTBranch } // branch node, not leaf node
func (s *nodeS[T]) hasData() bool       { return s.nType&NTData != 0 }        // has data?
func (s *nodeS[T]) isEmpty() bool       { return s.nType&NTData == 0 }        // no data?
func (s *nodeS[T]) Modified() bool      { return s.nType&NTModified != 0 }    // modification state
func (s *nodeS[T]) Description() string { return s.description }              // description field
func (s *nodeS[T]) Comment() string     { return s.comment }                  // comment field
func (s *nodeS[T]) Tag() any            { return s.tag }                      // tag field
func (s *nodeS[T]) Key() string         { return s.pathS }                    // key field is the full path of this node
func (s *nodeS[T]) KeyPiece() string    { return string(s.path) }             // key piece field for this node
func (s *nodeS[T]) IsLeaf() bool        { return s.nType&NTMask == NTLeaf }   // leaf node?
func (s *nodeS[T]) IsBranch() bool      { return s.nType&NTMask == NTBranch } // branch node?
func (s *nodeS[T]) HasData() bool       { return s.nType&NTData != 0 }        //nolint:revive //data field is valid?
func (s *nodeS[T]) Empty() bool         { return s.nType&NTData == 0 }        //nolint:revive //data field is empty?

func (s *nodeS[T]) readLockFor(cb func(*nodeS[T])) { //nolint:revive
	if s.rw != nil {
		s.rw.RLock()
		cb(s)
		s.rw.RUnlock()
	} else {
		cb(s)
	}
}

func (s *nodeS[T]) lockFor(cb func(*nodeS[T])) { //nolint:revive
	if s.rw != nil {
		s.rw.Lock()
		cb(s)
		s.rw.Unlock()
	} else {
		cb(s)
	}
}

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
	// if !s.isBranch() {
	s.readLockFor(func(s *nodeS[T]) {
		data = s.data
	})
	// }
	return
}

// SetData sets the Data field of a node.
func (s *nodeS[T]) SetData(data T) {
	// if !s.isBranch() {
	s.lockFor(func(s *nodeS[T]) {
		s.data = data
		s.nType |= NTData
	})
	// }
}

func (s *nodeS[T]) SetTTL(duration time.Duration, trie Trie[T], cb OnTTLRinging[T]) {
	trie.SetTTLFast(s, duration, cb)
}

// SetEmpty clear the Data field.
//
// Internally, SetEmpty sets Data field to zero value, and Tag field to nil, since v1.2.7+.
func (s *nodeS[T]) SetEmpty() {
	// if !s.isBranch() {
	// s.nType &= ^NTData
	s.lockFor(func(s *nodeS[T]) {
		var t T
		s.data = t
		s.tag = nil
	})
	// }
}

// SetComment sets the Description and Comment field.
func (s *nodeS[T]) SetComment(desc, comment string) { //nolint:revive
	// if s.isBranch() {
	// 	return
	// }
	s.description, s.comment = desc, comment
}

// SetTag sets the Tag field.
//
// You may save any value into a Tag field.
func (s *nodeS[T]) SetTag(tag any) { //nolint:revive
	// if s.isBranch() {
	// 	return
	// }
	s.tag = tag
}

func (s *nodeS[T]) StartsWith(ch rune) bool { //nolint:revive
	if len(s.path) == 0 {
		return false
	}
	return s.path[0] == ch
}

func (s *nodeS[T]) EndsWith(ch rune) bool { //nolint:revive
	kl := len(s.path)
	if kl == 0 {
		return false
	}
	return s.path[kl-1] == ch
}

func (s *nodeS[T]) endsWith(ch rune) bool { //nolint:revive
	kl := len(s.path)
	if kl == 0 {
		return false
	}
	return s.path[kl-1] == ch
}

// func (s *nodeS[T]) endsWithLite(ch rune) bool { //nolint:revive,unused
// 	return s.path[len(s.path)-1] == ch
// }

func (s *nodeS[T]) remove(item *nodeS[T]) (removed bool) { //nolint:revive
	if item == nil {
		return
	}
	// remove a child
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

func (s *nodeS[T]) insertInternal(word []rune, fullPath string, data T, trie *trieS[T], cb OnSetEx[T]) (node *nodeS[T], oldData any) {
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
		// eg2: insert '/app/:client/tokens' into '/app/:client/tokens/:token',
	}

	if cpl < wordLen {
		// eg: insert 'apple' into 'app'
		if cpl > 0 {
			word = word[cpl:] //nolint:revive
		}
		matched, child := base.matchChildren(word)
		if matched {
			var n Node[T]
			n, oldData = child.insert(word, fullPath, data, trie, cb)
			if nn, ok := n.(*nodeS[T]); ok {
				node = nn
				if cb != nil {
					cb(fullPath, oldData, node, trie)
				}
			}
		} else {
			node = base.insertAsLeaf(word, fullPath, data)
			if cb != nil {
				cb(fullPath, nil, node, trie)
			}
		}
	} else {
		// hit this node,
		base.nType |= NTData
		node, oldData, base.data = base, base.data, data
		if cb != nil {
			cb(fullPath, oldData, node, trie)
		}
	}
	return
}

func (s *nodeS[T]) split(pos int, word []rune) (newNode *nodeS[T]) {
	tip("[store/radix] split original path %q by word %q at pos %d", string(s.path), string(word), pos)

	// origPath, origPathS := s.path, s.pathS
	// if assertEnabled {
	// 	defer func() {
	// 		origPath, origPathS = []rune(origPathS), string(origPath)
	// 	}()
	// }

	d := len(s.path) - pos
	_ = word

	newNode = &nodeS[T]{
		path:        s.path[pos:],
		pathS:       s.pathS, // [pos:], // s.pathS[len(s.pathS)-d:],
		children:    s.children,
		data:        s.data,
		description: s.description,
		comment:     s.comment,
		tag:         s.tag,
		nType:       s.nType,
	}
	assert(strings.HasSuffix(newNode.pathS, string(newNode.path)), "newNode: pathS should end with path")

	s.path = s.path[:pos]
	s.pathS = s.pathS[:len(s.pathS)-d] // s.pathS[:pos] //
	s.children = []*nodeS[T]{newNode}
	s.nType = NTBranch
	s.description = ""
	s.comment = ""
	s.tag = nil
	var t T
	s.data = t
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
		// not a bug, when we need to compare the given word
		// with each of children, just the first char need
		// to be tested, since only have one child will
		// own the testing prefix, or only one child
		// have the part of the testing prefix.
		if child.path[0] == word[0] {
			matched = true
			break
		}
	}
	return
}

func extractor(from, to int, src, piece []rune, delimiter rune) (ret string, pos int, end bool) {
	j := 0
	for i := from; i < to; j++ {
		if src[i] == delimiter {
			ret, pos = string(piece[:j]), i
			break
		}
		piece[j] = src[i]
		i++
		if i == to {
			ret, pos, end = string(piece[:j+1]), i, true
		}
	}
	return
}

// search node by dotted path key with RecursiveMode.
func (s *nodeS[T]) search(mctx *matchCtx, word []rune, lastRuneIsDelimiter bool, parentNode *nodeS[T], kvpair KVPair) (matched, partialMatched bool, child, parent *nodeS[T]) { //nolint:revive
	matched, partialMatched, child, parent = s.matchR(mctx, word, lastRuneIsDelimiter, parentNode, kvpair)
	return
}

// matchR matches a path by walking child nodes recursively.
func (s *nodeS[T]) matchR(mctx *matchCtx, word []rune, lastRuneIsDelimiter bool, parentNode *nodeS[T], kvpair KVPair) (matched, partialMatched bool, child, parent *nodeS[T]) { //nolint:revive
	wl, l := len(word), len(s.path)
	if wl == 0 {
		return true, false, s, parentNode
	}

	// dm: delimiter just matched?
	// base: the working node ptr
	base, srcMatchedL, dstMatchedL, minL, maxL := s, 0, 0, min(l, wl), max(l, wl)
masterLoop:
	for ; srcMatchedL < minL; srcMatchedL++ {
		ch := base.path[srcMatchedL]
		if ch1 := word[srcMatchedL]; ch == ch1 {
			lastRuneIsDelimiter = ch == mctx.delimiter
			continue // first comparing loop, assume the index to base.path and word are both identical.
		}

		dstMatchedL = srcMatchedL // sync the index now

		// if partial matched,
		if srcMatchedL < l {
			if srcMatchedL < wl {
				var id, val string
				var srcEnd, dstEnd bool
			retriever:
				if lastRuneIsDelimiter {
					// matching "/*filepath"
					if ch == '*' {
						piece := make([]rune, maxL)
						id, srcMatchedL, srcEnd = extractor(srcMatchedL+1, l, base.path, piece, mctx.delimiter)
						if !srcEnd {
							logz.Warn("[matchR] invalid wildcard matching rule, it can only at end of the rule", "id", id, "srcMatchedL", srcMatchedL)
						}
						val, dstMatchedL, dstEnd = string(word[dstMatchedL:]), wl, true
						if kvpair != nil {
							kvpair[id] = val
						}
						logz.Verbose("[matchR] ident matched", "ident", id, "val", val, "key-path", base.pathS, "matching", string(word))
						// break masterLoop
						matched, child, parent = true, base, parentNode
						return
					}
					// matching "/:id/"
					if ch == ':' {
						// delimiter+':'+ident? | eg, matching source word "/hello/bob" on a trie-path pattern "/hello/:name"

						piece := make([]rune, maxL)
						id, srcMatchedL, srcEnd = extractor(srcMatchedL+1, l, base.path, piece, mctx.delimiter)
						val, dstMatchedL, dstEnd = extractor(dstMatchedL, wl, word, piece, mctx.delimiter)
						// logz.Verbose("[matchR] ident matched", "ident", id, "val", val, "key-path", base.pathS, "matching", string(word))
						// _, _, _, _, _ = dm, srcEnd, dstEnd, srcMatchedL, dstMatchedL

						if kvpair != nil {
							kvpair[id] = val
						}
						if srcEnd { // s.path matched to end.
							if dstEnd { // word matched to end.
								matched, child, parent = true, base, parentNode
								return
							}
							// not matched, break and match the rest part by looping base.children
							break masterLoop
						} else if !dstEnd {
							// sub-comparing loop here,
							for ; srcMatchedL < l && dstMatchedL < wl; srcMatchedL++ {
								ch = base.path[srcMatchedL]
								if ch1 := word[dstMatchedL]; ch == ch1 {
									lastRuneIsDelimiter = ch == mctx.delimiter
									dstMatchedL++
									continue
								}
								break // not matched, break to return false, false, nil, nil
							}
							if srcMatchedL == l && dstMatchedL == wl {
								matched, child, parent = true, base, parentNode
								break masterLoop // matched ok
							}
							if srcMatchedL == l && dstMatchedL < wl && len(base.children) > 0 {
								// get into and matchR with children nodes
								partialMatched = true
								break masterLoop
							}
							goto retriever // for the next id+val
						}
					}
				}
				// NOT matched and shall stop.
				// eg: matching 'apple' on 'apk'
				return false, false, nil, nil
			}
			// matched.
			// eg: matching 'app' on 'apple', or 'apk' on 'apple'
			return true, false, base, parentNode
		}
	}

	if srcMatchedL == l-1 && base.path[srcMatchedL] == mctx.delimiter {
		matched, child, parent = true, base, parentNode
	} else if minL < l && srcMatchedL == minL {
		partialMatched, child, parent = true, base, parentNode
	} else if minL >= l && srcMatchedL == minL && minL > 0 && srcMatchedL > 0 && !partialMatched {
		matched, child, parent = true, base, parentNode
	}
	if minL < wl {
		if len(base.children) == 0 {
			matched, partialMatched = false, true
			return
		}

		// restPart := word[minL:]
		if dstMatchedL == 0 {
			dstMatchedL = minL
		}
		restPart := word[dstMatchedL:]
		for _, child = range base.children {
			matched, partialMatched, child, parent = child.matchR(mctx, restPart, lastRuneIsDelimiter, s, kvpair)
			if matched || partialMatched {
				return
			}
		}
	}
	return
}

// matchL matches a path by iterating child nodes in a for-loop.

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

		s.readLockFor(func(n *nodeS[T]) {
			if s.hasData() {
				_, _ = sb.WriteString(" ")
				_, _ = sb.WriteString(s.pathS)
				_, _ = sb.WriteString(" => ")
				_, _ = sb.WriteString(ColorToDim(fmt.Sprint(s.Data())))
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
		})
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

	data := evendeep.MakeClone(s.Data())
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
