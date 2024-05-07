/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golistic/console"
)

func main() {
	var theme string

	flag.StringVar(&theme, "theme", "ascii", "Theme to use")
	flag.Parse()

	for _, what := range flag.Args() {
		switch what {
		case "toggle":
			toggle(theme)
		case "selection":
			selection(theme)
		}
	}
}

func toggle(theme string) {

	options := []string{"Absolutely!", "No.."}
	values := []bool{true, false}

	tg, err := console.NewToggle("Continue?", options, values)
	if err != nil {
		log.Fatal(err)
	}

	tg.SetTheme(theme)

	if err := tg.Render(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %v\n", tg.Label(), tg.Selected())
}

func selection(theme string) {

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
		log.Fatal(err)
	}

	fmt.Println("Selected:", s.Selected())
}
