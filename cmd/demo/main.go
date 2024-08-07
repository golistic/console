/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golistic/console"
)

func main() {
	var theme string

	flag.StringVar(&theme, "theme", "ascii", "Theme to use")
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage: demo [-theme=<theme>] [toggle|selection]")
	}

	for _, what := range flag.Args() {
		var err error
		switch what {
		case "toggle":
			err = toggle(theme)
		case "selection":
			err = selection(theme)
		}

		switch {
		case errors.Is(err, console.ErrCancelled):
			fmt.Println("Cancelled")
		case err != nil:
			fmt.Println("Error:", err)
		}
	}
}

func toggle(theme string) error {

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

func selection(theme string) error {

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
