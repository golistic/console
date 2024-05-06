/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import "fmt"

type Direction int

const (
	directionNone Direction = iota
	directionUp   Direction = iota
	directionDown
)

var cursorUp = string([]byte{27, 91, 65})
var cursorDown = string([]byte{27, 91, 66})

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}
