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
	directionLeft
	directionRight
)

var cursorUp = string([]byte{27, 91, 'A'})
var cursorDown = string([]byte{27, 91, 'B'})
var cursorRight = string([]byte{27, 91, 'C'})
var cursorLeft = string([]byte{27, 91, 'D'})

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}
