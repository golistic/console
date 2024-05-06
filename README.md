console - Tools and widgets for an interactive console
======================================================

Copyright (c) 2024, Geert JM Vanderkelen

Package `console` offers a collection of tools and widgets for creating an interactive
command line interface (CLI).

Quick Start
-----------

Show a list and use Up/Down keys to select an option:

```go
package main

import (
	"fmt"
	"log"

	"github.com/golistic/console"
)

func main() {

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

	if err := s.RenderWithTheme("ascii"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Selected:", s.Selected())
}
```

The above will adapt by default to the height of the terminal, but you can use `s.SetShowing(5)` to only show
5 options at a time.

Supported themes:
- `ascii`: simple, and works everywhere
- `nerdfont`: you need to use [Nerd Font in your terminal][1]


License
-------

Distributed under the MIT license. See `LICENSE.txt` for more information.

[1]: https://www.nerdfonts.com