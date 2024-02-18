package radix

import (
	"testing"
)

func TestGetpc(t *testing.T) {
	assert(t.Skipped(), "testing skipped")

	pc := getpc(2, 0, 1) // this
	src := getpcsource(pc)
	const mesg = "123"
	t.Logf("[ASSERT] condition failure | %s | %s:%d %s\n", mesg, src.File, src.Line, src.Function)

	str := stack(2, 0)
	t.Logf("%v", str)

	str = stack(0, 3)
	t.Logf("%v", str)
}
