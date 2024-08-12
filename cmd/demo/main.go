/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golistic/console"
)

func main() {
	var themeArg string

	flag.StringVar(&themeArg, "theme", "ascii", "Theme to use")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: demo [-theme=<theme>] [toggle|selection]")
	}

	theme := console.ThemeAscii
	if themeArg != "" {
		if t, ok := console.ThemeLookup(themeArg); !ok {
			fmt.Println("Theme must be one of: ascii, nerdfont, color01, inverted")
			os.Exit(1)
		} else {
			theme = t
		}
	}

	for _, what := range flag.Args() {
		var err error
		switch what {
		case "toggle":
			err = toggle(theme)
		case "selection":
			err = selection(theme)
		case "form":
			err = useForm(theme)
		}

		switch {
		case errors.Is(err, console.ErrAborted):
			fmt.Println("Cancelled")
		case err != nil:
			fmt.Println("Error:", err)
		}
	}
}

func toggle(theme console.Theme) error {

	options := []string{"Absolutely!", "No.."}
	values := []bool{true, false}

	tg, err := console.NewToggle("Continue?", options, values)
	if err != nil {
		log.Fatal(err)
	}

	tg.SetTheme(theme)

	if err := tg.Render(); err != nil {
		return err
	}

	fmt.Printf("%s %v\n", tg.Label(), tg.Selected())
	return nil
}

func selection(theme console.Theme) error {

	var options []string
	var values []int

	for i := range 40 {
		options = append(options, fmt.Sprintf("Option %02d", i))
	}

	for i := range 40 {
		values = append(values, i)
	}

	s, err := console.NewSelection(options, values)
	if err != nil {
		log.Fatal(err)
	}
	s.SetShowing(6)

	if err := s.RenderWithTheme(theme); err != nil {
		return err
	}

	fmt.Println("Selected:", s.Selected())
	return nil
}

func useForm(theme console.Theme) error {

	programmingLanugages := console.SelectProps{
		Options: []string{"Go", "Python", "TypeScript"},
		Values:  []any{"Go", "Python", "Typescript"},
	}

	scanner := func(value any, destination any) error {

		switch dest := destination.(type) {
		case *string:
			if v, ok := value.(string); !ok {
				return fmt.Errorf("expected string value")
			} else {
				*dest = v
			}
		case *int:
			switch v := value.(type) {
			case int:
				*dest = v
			case string:
				if n, err := strconv.Atoi(v); err != nil {
					return fmt.Errorf("expected int value")
				} else {
					*dest = n
				}
			default:
				return fmt.Errorf("expected string or int value")
			}
		default:
			return fmt.Errorf("unsupported type (was %T)", destination)
		}

		return nil
	}

	form := console.NewFormWithScanner(scanner).SetTheme(theme)

	var name string
	var favLang string
	var yearExp int

	form.AddElements(
		console.NewFormInput("name", "Your name", &name).AddValidator(func(value any) error {
			if name, ok := value.(string); !ok || len(name) == 0 {
				return errors.New("name is required")
			}
			return nil
		}),
		console.NewFormSelect("favLang", "Favorite language", &favLang, programmingLanugages),
		console.NewFormInput("yearExp", "Years Experience", &yearExp),
	)

	if err := form.Execute(); err != nil {
		return err
	}

	form.Clear()

	fmt.Printf("Hi %s! Your favorite language is %s and you have %d year(s) experience.\n",
		name, favLang, yearExp)
	return nil
}
