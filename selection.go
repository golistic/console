/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type selectionTheme struct {
	Unselected string
	Selected   string
}

var selectionThemes = map[string]selectionTheme{
	"nerdfont": {
		Unselected: "\uEBB5 %s",
		Selected:   "\u001B[32m\uF058 \u001B[0m%s",
	},
	"color01": {
		Unselected: "\u001B[47;30m%s\u001B[0m",   // BG:LightGrey FG:Black
		Selected:   "\u001B[1;42;30m%s\u001B[0m", // BG:Green FG:White
	},
	"ascii": {
		Unselected: "   %s",
		Selected:   "> %s",
	},
}

func NewSelection[V ~[]E, E any](options []string, values V) (*Selection[E], error) {

	if len(options) != len(values) {
		return nil, fmt.Errorf("number of options and values does not match")
	}

	s := &Selection[E]{
		values:  values,
		options: options,
		theme:   selectionThemes[defaultTheme],
	}

	s.SetShowing(0) // adapt to terminal height

	return s, nil
}

// Selection represents a selectable list of options with corresponding values.
// A user can use the Up- and Down-cursor keys to select an option, and push Enter
// to confirm the selection.
//
// By default, the `ascii` theme is used, but a Nerd Font theme `nerdfont` is also
// available.
type Selection[E any] struct {
	options []string
	values  []E

	showing int
	pointer int
	start   int
	end     int

	selectedOption E

	theme selectionTheme
}

func (s *Selection[E]) SetTheme(name string) {

	s.theme = selectionThemes[name]
}

// Selected returns the currently selected option from the Selection.
func (s *Selection[E]) Selected() E {

	return s.selectedOption
}

// SetShowing sets the number of options to be shown in the selection.
// If n is less than 1, it sets the number of options to the terminal height minus 3.
// Otherwise, it sets the number of options to n.
func (s *Selection[E]) SetShowing(n int) {

	if n < 1 {
		_, s.showing = s.terminalSize()
		s.showing -= 3
	} else {
		s.showing = n
	}
}

func (s *Selection[E]) terminalSize() (int, int) {

	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		return 80, 24
	}

	return width, height
}

// RenderWithTheme renders the Selection with the specified theme. If the theme with the given
// name does not exist, the default theme of the Selection is used.
func (s *Selection[E]) RenderWithTheme(themeName string) error {

	theme, ok := selectionThemes[themeName]
	if !ok {
		theme = s.theme
	}

	return s.render(theme)
}

// Render renders the Selection.
func (s *Selection[E]) Render() error {
	return s.render(s.theme)
}

func (s *Selection[E]) render(theme selectionTheme) error {

	s.pointer = 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)

		for i := 0; i < s.showing; i++ {
			fmt.Print(cursorUp)
			fmt.Print("\r\033[2K")
		}

		fmt.Print("\r\033[2K")
		showCursor()
	}()

	hideCursor()

	if len(s.options) < s.showing-1 {
		s.showing = len(s.options)
	}

	s.renderOptions(theme, s.options, directionNone)

	var done bool
	for {
		if done {
			break
		}
		b := make([]byte, 3)
		if _, err := os.Stdin.Read(b); err != nil {
			return fmt.Errorf("reading input (%w)", err)
		}

		switch {
		case b[0] == 10 || b[0] == 13:
			s.selectedOption = s.values[s.pointer]
			done = true
		case b[0] == 27 && b[1] == 91:
			var direction Direction

			switch b[2] {
			case cursorUp[2]:
				direction = directionUp
			case cursorDown[2]:
				direction = directionDown
			default:
				direction = directionNone
			}

			if direction != directionNone {
				for i := 0; i < s.showing; i++ {
					fmt.Print(cursorUp)
				}
				s.renderOptions(theme, s.options, direction)
			}
		case b[0] == 3 || b[0] == 27:
			return ErrCancelled
		}
	}

	return nil
}

func (s *Selection[E]) renderOptions(theme selectionTheme, options []string, direction Direction) {

	lenOpts := len(options)

	if direction == directionUp && s.pointer > 0 {
		s.pointer--
	} else if direction == directionDown && s.pointer < len(s.options)-1 {
		s.pointer++
	}

	if direction == directionUp {
		if s.pointer < s.start {
			if s.pointer >= lenOpts-s.showing {
				s.start = lenOpts - s.showing - 1
			} else {
				s.start = s.pointer
			}
		}
	} else if direction == directionDown && s.pointer >= s.showing {
		if s.pointer < lenOpts {
			s.start = s.pointer - s.showing + 1
		}
	}

	s.end = s.start + s.showing

	// handle when at end of options
	if s.end > lenOpts {
		s.end = lenOpts
		s.start = s.end - s.showing
		if s.start < 0 {
			s.start = 0
		}
	}

	for i := s.start; i < s.end; i++ {

		if i == s.pointer {
			fmt.Printf("\r\033[2K %s\n",
				fmt.Sprintf(theme.Selected, options[i]))
		} else {
			fmt.Printf("\r\033[2K %s\n",
				fmt.Sprintf(theme.Unselected, options[i]))
		}
	}
}
