/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import "strings"

const defaultTheme = "ascii"

type Theme string

const (
	ThemeNerdFont Theme = "nerdfont"
	ThemeAscii    Theme = "ascii"
	ThemeColor01  Theme = "color01"
	ThemeInverted Theme = "inverted"
)

var AllTheme = []Theme{
	ThemeNerdFont,
	ThemeAscii,
	ThemeColor01,
	ThemeInverted,
}

func ThemeLookup(name string) (Theme, bool) {
	name = strings.ToLower(name)
	for _, t := range AllTheme {
		if string(t) == name {
			return t, true
		}
	}

	return "", false
}
