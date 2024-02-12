//go:build !delve && !verbose

package radix

func tip(msg ...any) {} //nolint:revive

const tipEnabled = false //nolint:unused
