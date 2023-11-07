//go:build !delve
// +build !delve

package radix

const assertEnabled = false

var assertAlwaysStop = false

func assert(cond bool, msg ...any) {}
