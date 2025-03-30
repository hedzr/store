package radix

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/hedzr/errors.v3"

	logz "github.com/hedzr/logg/slog"
)

// NewTrie returns a Trie-tree instance.
func NewTrie[T any]() *trieS[T] {
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

// NewTrieBy returns a Trie-tree instance.
func NewTrieBy[T any](delimiter rune) *trieS[T] {
	return &trieS[T]{root: &nodeS[T]{}, delimiter: delimiter}
}

var _ Trie[any] = (*trieS[any])(nil) // assertion helper

func newTrie[T any]() *trieS[T] { //nolint:revive
	return &trieS[T]{root: &nodeS[T]{}, delimiter: dotChar}
}

type trieS[T any] struct {
	root       *nodeS[T]
	prefix     string
	delimiter  rune
	ttlpresent atomic.Uint32
	ttls       *TTL[T]
}

type TTL[T any] struct {
	treevec []*trieS[T]
	mu      sync.RWMutex
	cancel  context.CancelFunc // exit signal
	done    <-chan struct{}
	adder   chan ttljobS[T]

	// in the future, we might try timing-wheel way.

	weeks   *wheelS // by weeks, years, ...
	years   *wheelS // by years, months, days, weeks
	days    *wheelS // by days, hours, minutes, seconds
	seconds *wheelS // by seconds, ms, us, ns
}

type wheelS struct {
	wheel    map[int]*wheelS
	mu       sync.RWMutex
	onAction wheelAction
}

type wheelAction func(ctx context.Context, w *wheelS)

type ttljobS[T any] struct {
	node     *nodeS[T]
	duration time.Duration
	action   OnTTLRinging[T]
}

type OnTTLRinging[T any] func(s *TTL[T], nd Node[T])

func newttls[T any](t *trieS[T]) *TTL[T] {
	ctx, cancel := context.WithCancel(context.Background())
	s := &TTL[T]{
		treevec: []*trieS[T]{t},
		cancel:  cancel,
		done:    ctx.Done(),
		adder:   make(chan ttljobS[T]),
	}
	go s.run()
	return s
}

func (s *TTL[T]) dupS() *TTL[T] {
	ctx, cancel := context.WithCancel(context.Background())
	n := &TTL[T]{
		treevec: []*trieS[T]{s.treevec[0]},
		cancel:  cancel,
		done:    ctx.Done(),
		adder:   make(chan ttljobS[T]),
	}
	go n.run()
	return n
}

func (s *TTL[T]) Close() {
	s.cancel()
}

func (s *TTL[T]) Tree() Trie[T] { return s.treevec[0] }

func (s *TTL[T]) Add(nd *nodeS[T], duration time.Duration, action OnTTLRinging[T]) {
	if nd.rw == nil {
		nd.rw = &sync.RWMutex{}
	}
	s.adder <- ttljobS[T]{node: nd, duration: duration, action: action}
}

func (s *TTL[T]) run() {
	adder := func(job ttljobS[T]) {
		if job.duration < 200*time.Nanosecond {
			if job.duration == 0 {
				// reset a job (might be planned in the future)
			}
			return
		}

		// add a new job
		timer := time.NewTimer(job.duration)
		go func(timer *time.Timer, job ttljobS[T]) {
			defer timer.Stop()
			// defer func() { job.node.rw = nil }()
			for {
				select {
				case <-timer.C:
					if job.action != nil {
						job.action(s, job.node)
					}
					if job.node.isBranch() {
						s.treevec[0].Remove(job.node.pathS)
					} else {
						job.node.SetEmpty()
					}
					return
				case <-s.done:
					return
				}
			}
		}(timer, job)
	}

	for {
		select {
		case <-s.done:
			return
		case job := <-s.adder:
			adder(job)
		}
	}
}

func (s *trieS[T]) Close() {
	if s.ttlpresent.CompareAndSwap(1, 0) {
		if s.ttls != nil {
			s.ttls.Close()
		}
	}
}

// SetTTL sets a ttl timeout for a branch or a leaf node.
//
// At ttl arrived, the leaf node value will be cleared.
// For a branch node, it will be dropped.
//
// Once you're using SetTTL, don't forget call Close().
// For example:
//
//	conf := newTrieTree()
//	defer conf.Close()
//
//	path := "app.verbose"
//	conf.SetTTL(path, 200*time.Millisecond, func(ctx context.Context, s *TTL[any], nd *Node[any]) {
//		t.Logf("%q cleared", path)
//	})
//
// **[Pre-API]**
//
// SetTTL is a prerelease API since v1.2.5, it's mutable in the
// several future releases recently.
//
// The returned `state`: 0 assumed no error.
func (s *trieS[T]) SetTTL(path string, ttl time.Duration, cb OnTTLRinging[T]) (state int) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	found := node != nil && !partialMatched // && !node.isBranch()
	state = -1
	if found {
		state = 0
		if s.ttlpresent.CompareAndSwap(0, 1) {
			s.ttls = newttls(s)
		}
		s.ttls.Add(node, ttl, cb)
	}
	return
}

func (s *trieS[T]) SetTTLFast(node Node[T], ttl time.Duration, cb OnTTLRinging[T]) (state int) {
	if nd, ok := node.(*nodeS[T]); ok {
		if s.ttlpresent.CompareAndSwap(0, 1) {
			s.ttls = newttls(s)
		}
		s.ttls.Add(nd, ttl, cb)
	} else {
		state = -1
	}
	return
}

//

// dupS for duplicating itself. see also Dup, WithPrefix, WithPrefix & WithPrefixReplaced, withPrefixReplacedImpl.
func (s *trieS[T]) dupS(root *nodeS[T], prefix string) (newTrie *trieS[T]) { //nolint:revive
	newTrie = &trieS[T]{
		root:      root,
		prefix:    prefix,
		delimiter: s.delimiter,
	}
	if s.ttlpresent.Load() > 0 {
		newTrie.ttls = s.ttls.dupS()
		newTrie.ttlpresent.Add(1)
	}
	return
}

func (s *trieS[T]) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString("Trie{\"")
	_, _ = sb.WriteString(s.prefix)
	_, _ = sb.WriteString("\", delimiter:'")
	_, _ = sb.WriteRune(s.delimiter)
	_, _ = sb.WriteString("'}")
	return sb.String()
}

func (s *trieS[T]) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	_, _ = sb.WriteString("\"trie\":{\"prefix\":")
	_, _ = sb.WriteString(strconv.Quote(s.prefix))
	_, _ = sb.WriteString(",\"delimiter\":\"")
	_, _ = sb.WriteRune(s.delimiter)
	_, _ = sb.WriteString("\"}")
	return []byte(sb.String()), nil
}

//

//

//

const (
	dotChar                 rune = '.'
	initialPrefixBufferSize      = 64
	maxPrefixBufferSize          = 64*1024*1024 - initialPrefixBufferSize
)

var prefixJointPool = sync.Pool{New: func() any {
	return bytes.NewBuffer(make([]byte, 0, initialPrefixBufferSize))
}}

func (s *trieS[T]) poolGet() *bytes.Buffer {
	return prefixJointPool.Get().(*bytes.Buffer) //nolint:revive
}

func (s *trieS[T]) deferPoolGet(bb *bytes.Buffer) {
	bb.Reset()
	if bb.Cap() < maxPrefixBufferSize {
		prefixJointPool.Put(bb)
	}
}

func (s *trieS[T]) UsePool(fn func(bb *bytes.Buffer, bytes int)) string {
	i, bb := 0, s.poolGet()
	defer s.deferPoolGet(bb)
	fn(bb, i)
	return bb.String()
}

func (s *trieS[T]) join1(pre string, args ...string) (ret string) {
	if pre == "" {
		return s.Join(args...)
	}

	if len(args) == 0 {
		return pre
	}

	i, bb := 0, s.poolGet()
	defer s.deferPoolGet(bb)

	_, _ = bb.WriteString(pre)

	for _, it := range args {
		if it != "" {
			_ = bb.WriteByte(byte(s.delimiter))
			_, _ = bb.WriteString(it)
			i++
		}
	}
	return bb.String()
}

func (s *trieS[T]) Join(args ...string) (ret string) {
	switch len(args) {
	case 0:
		return
	case 1:
		return args[0]
	}

	if args[0] == "" {
		return s.Join(args[1:]...)
	}

	return s.UsePool(func(bb *bytes.Buffer, bytes int) {
		for _, it := range args {
			if it != "" {
				if bytes > 0 {
					_ = bb.WriteByte(byte(s.delimiter))
				}
				_, _ = bb.WriteString(it)
				bytes++ //nolint:revive
			}
		}
	})

	// i, bb := 0, s.poolGet()
	// defer s.deferPoolGet(bb)
	//
	// for _, it := range args {
	// 	if it != "" {
	// 		if i > 0 {
	// 			_ = bb.WriteByte(byte(s.delimiter))
	// 		}
	// 		_, _ = bb.WriteString(it)
	// 		i++
	// 	}
	// }
	// return bb.String()
}

func (s *trieS[T]) Insert(path string, data T) (oldData any) { //nolint:revive
	_, oldData = s.Set(path, data)
	return
}

// Set sets the Data field into a node specified by path.
//
// If the given path cannot be found, a new node will be created at that
// location so that the new data value can be set into it.
func (s *trieS[T]) Set(path string, data T) (node Node[T], oldData any) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	return s.root.insert([]rune(path), path, data, s, nil)
}

type OnSetEx[T any] func(path string, oldData any, node Node[T], trie Trie[T])

func (s *trieS[T]) SetEx(path string, data T, cb OnSetEx[T]) (oldData any) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	_, oldData = s.root.insert([]rune(path), path, data, s, cb)
	return
}

// SetNode sets the all node fields at once.
func (s *trieS[T]) SetNode(path string, data T, tag any, descriptionAndComments ...string) (ret Node[T], oldData any) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, old := s.root.insertInternal([]rune(path), path, data, s, nil)
	switch len(descriptionAndComments) {
	case 0:
	case 1:
		node.description = descriptionAndComments[0]
	case 2:
		node.description, node.comment = descriptionAndComments[0], descriptionAndComments[1]
	default:
		node.description = descriptionAndComments[0]
		node.comment = strings.Join(descriptionAndComments[1:], "\n")
	}
	ret, oldData = node, old
	return
}

// SetEmpty clear the Data field.
func (s *trieS[T]) SetEmpty(path string) (oldData any) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	var v T
	node, old := s.root.insertInternal([]rune(path), path, v, s, nil)
	node.SetEmpty()
	return old
}

func (s *trieS[T]) Update(path string, cb func(node Node[T], old any)) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	var v T
	node, old := s.root.insertInternal([]rune(path), path, v, s, nil)
	cb(node, old)
	return
}

// SetComment sets the Desc and Comment field of a node specified by path.
//
// Nothing happens if the given path cannot be found.
func (s *trieS[T]) SetComment(path, description, comment string) (ok bool) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	if ok = node != nil || partialMatched; ok {
		node.description, node.comment = description, comment
	}
	return
}

// SetTag sets the Tag field of a node specified by path.
//
// Nothing happens if the given path cannot be found.
func (s *trieS[T]) SetTag(path string, tag any) (ok bool) { //nolint:revive// set extra notable data bound to a key
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	if ok = node != nil || partialMatched; ok {
		node.tag = tag
	}
	return
}

// Merge a map at path point 'pathAt'
func (s *trieS[T]) Merge(pathAt string, data map[string]any) (err error) {
	// if s.prefix != "" {
	// 	pathAt = s.Join(s.prefix, pathAt) //nolint:revive
	// }
	err = s.withPrefixImpl(pathAt).loadMap(data)
	return
}

// Search checks the path if it exists. = StartsWith
//
// Only fully-matched node are considered as FOUND.
// Which means, if a path matched partial part, suppose matching
// `a.bcd.e` in a tree has `a.bcd.ends` node, it will be
// matched as partial-match state. But it is non-FOUND.
// And if a tree has 'a.bcd.e.nds' node, FOUND returns.
//
// Using Location to retrieve more info for seaching a path.
func (s *trieS[T]) Search(path string) (found bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

// Locate checks a path if it exists.
func (s *trieS[T]) Locate(path string, kvpair KVPair) (node *nodeS[T], branch, partialMatched, found bool) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched = s.search(path, kvpair)
	if node != nil {
		found = !partialMatched
		branch = node.isBranch()
	}
	// found, branch = node != nil && !partialMatched, safeIsBranch(node)
	return
}

func safeIsBranch[T any](node *nodeS[T]) bool { return node != nil && node.isBranch() }

// Has tests of a path exists.
//
// Delimiter shall fully match a path (typically is dot '.').
// That means, for an existed tree:
//
//	app.lite-mode
//	app.logging
//	app.logging.enabled
//	app.logging.file
//
// Has("app.logging") or Has("app.logging.enabled") will get a true.
//
// But Has("app.logging.en") returns false.
//
// And Has("app.l") must be false.
func (s *trieS[T]) Has(path string) (found bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	found = node != nil && !partialMatched // && !node.isBranch()
	return
}

// HasPart tests if a path exists even if partial matched.
//
// Using Location to retrieve more info for searching a path.
func (s *trieS[T]) HasPart(path string) (yes bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, nil)
	yes = node != nil || partialMatched
	if partialMatched && node != nil {
		yes = strings.HasPrefix(node.pathS, path)
	}
	return
}

// Remove deleting a path from this tree.
//
// The return boolean represents there was a node removed.
// If the path does not exist, it will return false.
func (s *trieS[T]) Remove(path string) (removed bool) { //nolint:revive
	_, _, removed = s.RemoveEx(path)
	return
}

// RemoveEx deleting a path and return more status than Remove.
func (s *trieS[T]) RemoveEx(path string) (nodeRemoved, nodeParent Node[T], removed bool) {
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, parent, partialMatched := s.search(path, nil)
	found := node != nil && !partialMatched // && !node.isBranch()
	if found {
		if parent != nil {
			removed = parent.remove(node)
			if removed {
				nodeRemoved, nodeParent = node, parent
			}
		} else {
			logz.Warn("if given path found and return node, its parent MUST NOT be nil", "node", node, "parent", parent)
		}
	}
	return
}

// MustGet is a simple Get without a checked found state.
//
// If nothing is found, zero data returned.
func (s *trieS[T]) MustGet(path string) (data T) {
	var branch, found bool
	data, branch, found, _ = s.Query(path, nil)
	if !found && !branch {
		data = *new(T)
	}
	return
}

// Get searches the given path and return its data field if found.
func (s *trieS[T]) Get(path string) (data T, found bool) {
	data, _, found, _ = s.Query(path, nil)
	return
}

// Query searches a path and returns the located info: 'found' boolean flag
// identify the path found or not; 'branch' flag identify the found node
// is a branch or a leaf; for a leaf node, 'data' return its Data field.
//
// If something is wrong, 'err' might collect the reason for why. But,
// it generally is errors.NotFound (errors.Code -5).
func (s *trieS[T]) Query(path string, kvpair KVPair) (data T, branch, found bool, err error) { //nolint:revive
	if s.prefix != "" {
		path = s.Join(s.prefix, path) //nolint:revive
	}
	node, _, partialMatched := s.search(path, kvpair)
	found = node != nil && !partialMatched
	if found {
		if node.isBranch() {
			branch = true
			if !node.endsWith(s.delimiter) {
				found = false
			}
		}
		if node.hasData() {
			node.lockFor(func(n *nodeS[T]) {
				data = node.data
			})
		}
	}
	// if !found {
	// 	err = errors.NotFound
	// }
	err = iif(found, error(nil), error(errors.NotFound))
	return
}

func (s *trieS[T]) search(word string, kvpair KVPair) (found, parent *nodeS[T], partialMatched bool) { //nolint:revive
	found = s.root
	mctx := getMatchCtx(word, s.delimiter)
	// stringtoslicerune needs two pass full-scanning for a string, but it have to be to do.
	if matched, pm, child, prnt := found.matchR(mctx, []rune(word), false, nil, kvpair); matched || pm {
		putBack(mctx)
		return child, prnt, pm
	}
	putBack(mctx)
	found = nil
	return
}

func getMatchCtx(word string, delimiter rune) *matchCtx {
	s := matchCtxPool.Get().(*matchCtx)
	s.fullPath, s.delimiter = word, delimiter
	return s
}
func putBack(mctx *matchCtx) { matchCtxPool.Put(mctx) }

var matchCtxPool = sync.Pool{New: func() any {
	return &matchCtx{}
}}

type matchCtx struct {
	fullPath  string
	delimiter rune
}

type KVPair map[string]string

// Delimiter returns the current delimiter in using.
//
// A trieS-tree is a radix-tree, so child nodeS may split at any
// posistion in characters of a full path.
//
// But Query or others APIs make the path is meaningful for
// the using delimiter, like dot char ('.'). This is why
// there is 'partialMatched' in Locate returning values.
//
// The delimiter can be changed at runtime. After it's changed,
// the next Locate, Query or any others will interpret the
// path with new delimiter.
//
// Since the delimiter character is part of a full path,
// we can re-interpret its meaning dynamically without extras
// costs.
//
// And the benefits are not only replacing the delimiter
// dynamically: splitting a path by a delimiter or joining
// the segements splitted into from path are both unnecessary, in
// a conventional designing and implementing mode.
func (s *trieS[T]) Delimiter() rune { return s.delimiter }

// SetDelimiter sets the delimiter rune.
//
// When extracting a node and its data, the delimiter character is
// the decisive factor.
//
// The Store's default delimiter is '.' (dot). But you can construct
// a Trie-tree with other char, such as path separator ('/').
//
// For example,
//
//	trie := newTrie[any]()
//	trie.Insert("/search", 1)
//	trie.Insert("/support", 2)
//	trie.Insert("/blog/:post/", 3)
//	trie.Insert("/about-us/team", 4)
//	trie.Insert("/contact", 5)
//	trie.Insert("/about-us/legal", 6)
//
//	trie.SetDelimiter('/')
//	data, err := trie.GetM("/about-us")
//	assert.True(reflect.DeepEqual(data, map[string]any{"legal": 6, "team": 4}))
//
// See also TestTrieS_Delimiter(),
func (s *trieS[T]) SetDelimiter(delimiter rune) { s.delimiter = delimiter }

func (s *trieS[T]) simpleEndsWith(str string, ch rune) bool { //nolint:revive
	if str != "" {
		runes := []rune(str)
		return runes[len(runes)-1] == ch
	}
	return false
}

// Dump collects the nodes in this tree and prints them for a
// debugging purpose.
// It returns the formatted string then you can print it to
// stdout or a file.
//
//	println(trie.Dump())
//
// The dump results is decorated with ANSI escaped sequences.
// So if you want a plain pure text, enable NoColor mode
// defined in hedzr/is/states package:
//
//	import "github.com/hedzr/is/states"
//
//	states.Env().SetNoColorMode(true)
//	println(trie.Dump())
//
// Or, [StatesEnvSetColorMode(true)] can also do that.
func (s *trieS[T]) Dump() string             { return s.root.dump(false) }   //nolint:revive
func (s *trieS[T]) dump(noColor bool) string { return s.root.dump(noColor) } //nolint:revive

// Dup or Clone makes an exact deep copy of this tree.
func (s *trieS[T]) Dup() (newTrie *trieS[T]) { //nolint:revive
	return s.dupS(s.root.Dup(), s.prefix)
}

// Walk navigates the whole tree (passing "" as 'path' param) or
// a subtree from a given path.
func (s *trieS[T]) Walk(path string, cb func(path, fragment string, node Node[T])) { //nolint:revive
	root := s.root
	if path != "" {
		node, parent, partialMatched := s.search(path, nil)
		if !partialMatched {
			root = parent
			if runes := []rune(path); runes[len(runes)-1] == s.delimiter {
				root = node
			}
		}
	}

	if root != nil {
		root.walk(0, cb)
	}
}
