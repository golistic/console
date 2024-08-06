/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"os"

	"golang.org/x/term"
)

// TerminalSize returns the visible dimensions of the given terminal.
// When size could not be retrieved, weight 80 and height 24 is returned.
func TerminalSize() (width, height int) {

	fd := int(os.Stdout.Fd())

	var err error
	width, height, err = term.GetSize(fd)
	if err != nil {
		return 80, 24
	}

	return width, height
}
