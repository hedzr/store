//go:build !delve && !verbose

package radix

func tip(msg ...any) {}

const tipEnabled = false
