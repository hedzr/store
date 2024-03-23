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
var ColorToDim = color.ToDim

// ColorToColor is a wrapper for [color.ToColor] to make a string around by specified [color.Color].
var ColorToColor = color.ToColor

const (
	FgLightGreen = color.FgLightGreen // light green
	FgGreen      = color.FgGreen      // green
)
