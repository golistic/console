/*
 * Copyright (c) 2025, Geert JM Vanderkelen
 */

package console

import "fmt"

func NewFormText(text string) *FormText {
	return &FormText{
		text:        text,
		formElement: &formElement{},
	}
}

type FormText struct {
	*formElement
	text string
}

var _ FormElementer = (*FormInput)(nil)

func (ft *FormText) do() error {

	fmt.Printf("%s\n", ft.text)
	return nil
}
