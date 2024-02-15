package store

import (
	"fmt"
	"strings"

	"gopkg.in/hedzr/errors.v3"
)

type Op uint32 // Op describes a set of file operations.

var opStrings = map[Op]string{
	OpCreate: "create",
	OpWrite:  "modify",
	OpRename: "rename",
	OpRemove: "remove",
	OpChmod:  "chmod",
	OpNone:   "none",
}

var opStringsRev = map[string]Op{
	"create": OpCreate,
	"new":    OpCreate,
	"modify": OpWrite,
	"write":  OpWrite,
	"rename": OpRename,
	"remove": OpRemove,
	"delete": OpRemove,
	"rm":     OpRemove,
	"chmod":  OpChmod,
	"none":   OpNone,
}

func (s *Op) UnmarshalText(text []byte) error {
	// panic("implement me")
	op, ok := opStringsRev[string(text)]
	if ok {
		*s = op
		return nil
	}
	return errors.New("bad/unknown string, can't unmarshal to Op")
}

func (s *Op) MarshalText() (text []byte, err error) {
	sz, ok := opStrings[*s]
	if ok {
		return []byte(strings.ToUpper(sz)), nil
	}
	return []byte(fmt.Sprintf("Op(%d)", s)), nil
}

func (s *Op) Marshal() []byte {
	return nil
}

// The operations fsnotify can trigger; see the documentation on [Watcher] for a
// full description, and check them with [Event.Has].
const (
	// OpCreate is a new pathname was created.
	OpCreate Op = 1 << iota

	// OpWrite the pathname was written to; this does *not* mean the write has finished,
	// and a write can be followed by more writes.
	OpWrite

	// OpRemove the path was removed; any watches on it will be removed. Some "remove"
	// operations may trigger a Rename if the file is actually moved (for
	// example "remove to trash" is often a rename).
	OpRemove

	// OpRename the path was renamed to something else; any watched on it will be
	// removed.
	OpRename

	// OpChmod file attributes were changed.
	//
	// It's generally not recommended to take action on this event, as it may
	// get triggered very frequently by some software. For example, Spotlight
	// indexing on macOS, anti-virus software, backup software, etc.
	OpChmod

	OpNone = 0
)
