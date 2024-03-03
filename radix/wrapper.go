package radix

import (
	"github.com/hedzr/is/states"
	"github.com/hedzr/is/term/color"
)

// StatesEnvSetColorMode is a setter for "--no-color"
var StatesEnvSetColorMode = func(b bool) {
	states.Env().SetNoColorMode(b)
}

// ColorToDim is a wrapper for [color.ToDim] to make a string around by dimmed foreground.
var ColorToDim = func(format string, args ...any) string {
	return color.ToDim(format, args...)
}

// ColorToColor is a wrapper for [color.ToColor] to make a string around by specified [color.Color].
var ColorToColor = func(clr color.Color, format string, args ...any) string {
	return color.ToColor(clr, format, args...)
}

const (
	FgLightGreen = color.FgLightGreen // light green
	FgGreen      = color.FgGreen      // green
)
