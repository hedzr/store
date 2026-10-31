package radix

import "fmt"

// RecursiveMode specifies how Must/GetXXX looks up a key
// for matching its parent nodes (Up) or children (Down)
// if not matched current node.
//
// By default, they don't (RecursiveNone).
//
// See also: [Trie.N], [Trie.R] and [Trie.BR]
type RecusiveMode int

const (
	RecusiveNone RecusiveMode = iota
	RecusiveUp                // match path-key recusively up for Must/GetXXX
	RecusiveDown              // match path-key recusively down for Must/GetXXX
)

func (e RecusiveMode) String() string {
	switch e {
	case RecusiveNone:
		return "None"
	case RecusiveUp:
		return "Up"
	case RecusiveDown:
		return "Down"
	default:
		return fmt.Sprintf("RrMode(#%d)", int(e))
	}
}
