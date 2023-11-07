//go:build delve || verbose
// +build delve verbose

package radix

import (
	"fmt"
	"os"
)

func tip(msg ...any) {
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
	fmt.Fprintf(os.Stderr, "[LOG] %s | %s:%d %s\n", mesg, src.File, src.Line, src.Function)
}

const tipEnabled = true
