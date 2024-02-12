//go:build !delve
// +build !delve

package radix

const assertEnabled = false //nolint:unused

var assertAlwaysStop = false //nolint:revive,unused

func assert(cond bool, msg ...any) {} //nolint:revive
