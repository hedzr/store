//go:build delve
// +build delve

package radix

import (
	"fmt"
	"os"
)

const assertEnabled = true

var assertAlwaysStop = true

func assert(cond bool, msg ...any) {
	if !cond {
		pc := getpc(2) // this
		src := getpcsource(pc)
		var mesg string
		if len(msg) > 0 {
			if format, ok := msg[0].(string); ok {
				mesg = fmt.Sprintf(format, msg[1:]...)
			} else {
				mesg = fmt.Sprint(msg...)
			}
		}
		fmt.Fprintf(os.Stderr, "[ASSERT] condition failure | %s | %s:%d %s\n", mesg, src.File, src.Line, src.Function)
		if assertAlwaysStop {
			panic("[ASSERT] condition failure")
		}
	}
}
