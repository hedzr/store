package store

import (
	"testing"

	"github.com/hedzr/is/stringtool"
)

func TestRandomStringPure(t *testing.T) {
	t.Log(stringtool.RandomStringPure(8))
}
