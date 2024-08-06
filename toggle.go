/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type toggleTheme struct {
	Unselected string
	Selected   string
	Gap        int
}

var toggleThemes = map[string]toggleTheme{
	"nerdfont": {
		Unselected: "\uEBB5 %s",
		Selected:   "\u001B[32m\uF058 \u001B[0m%s",
		Gap:        1,
	},
	"inverted": {
		Unselected: "%s",
		Selected:   "\u001B[7m%s\u001B[0m", // inverted
		Gap:        1,
	},
	"color01": {
		Unselected: "\u001B[47;30m%s\u001B[0m",   // BG:LightGrey FG:Black
		Selected:   "\u001B[1;42;30m%s\u001B[0m", // BG:Green FG:White
		Gap:        1,
	},
	"ascii": {
		Unselected: "  %s",
		Selected:   "> %s",
		Gap:        2,
	},
}

func NewToggle[V ~[]T, T any](label string, options []string, values V) (*Toggle[T], error) {

	if len(options) != 2 || len(options) != len(values) {
		return nil, fmt.Errorf("number of options and values must be 2")
	}

	toggle := &Toggle[T]{
		label:   label,
		options: options,
		values:  values,
	}

	toggle.SetTheme(defaultTheme)

	return toggle, nil
}

type Toggle[T any] struct {
	label   string
	options []string
	values  []T

	pointer        int
	selectedOption T

	theme toggleTheme
	gap   int
}

func (tg *Toggle[E]) SetTheme(name string) {

	theme, ok := toggleThemes[name]
	if !ok {
		theme = tg.theme
	}

	tg.theme = theme
}

// Selected returns the currently toggled option.
func (tg *Toggle[T]) Selected() T {

	return tg.selectedOption
}

func (tg *Toggle[T]) Label() string {
	return tg.label
}

// RenderWithTheme renders the Selection with the specified theme. If the theme with the given
// name does not exist, the default theme of the Selection is used.
func (tg *Toggle[T]) RenderWithTheme(themeName string) error {

	theme, ok := toggleThemes[themeName]
	if !ok {
		theme = tg.theme
	}

	return tg.render(theme)
}

// Render renders the Selection.
func (tg *Toggle[T]) Render() error {
	return tg.render(tg.theme)
}

func (tg *Toggle[T]) render(theme toggleTheme) error {

	tg.pointer = 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)

		fmt.Print("\r\033[2K")
		showCursor()
	}()

	hideCursor()
	tg.renderOptions(tg.theme, tg.options)

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
			tg.selectedOption = tg.values[tg.pointer]
			done = true
		case b[0] == 27 && b[1] == 91:
			switch b[2] {
			case cursorLeft[2]:
				tg.pointer = 0
			case cursorRight[2]:
				tg.pointer = 1
			}

			tg.renderOptions(theme, tg.options)
		case b[0] == 3 || b[0] == 27:
			return ErrCancelled
		}
	}

	return nil
}

func (tg *Toggle[T]) renderOptions(theme toggleTheme, options []string) {

	fmt.Printf("\r\033[2K%s ", tg.label)

	if tg.pointer == 0 {
		fmt.Printf("%s%s%s",
			fmt.Sprintf(theme.Selected, options[0]),
			strings.Repeat(" ", theme.Gap),
			fmt.Sprintf(theme.Unselected, options[1]))
	} else {
		fmt.Printf("%s%s%s",
			fmt.Sprintf(theme.Unselected, options[0]),
			strings.Repeat(" ", theme.Gap),
			fmt.Sprintf(theme.Selected, options[1]))
	}
}
