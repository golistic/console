/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"fmt"
)

type SelectProps struct {
	Options []string
	Values  []any
	Showing int
}

func NewFormSelect(name, label string, dest any, props SelectProps) *FormSelect {
	return &FormSelect{
		formElement: &formElement{
			name:  name,
			label: label,
			dest:  dest,
		},

		props: props,
	}
}

type FormSelect struct {
	*formElement

	props SelectProps
}

var _ FormElementer = (*FormSelect)(nil)

func (fs *FormSelect) setForm(form *Form) {

	fs.form = form
}

func (fs *FormSelect) do() error {

	selection, err := NewSelection(fs.props.Options, fs.props.Values)
	if err != nil {
		return err
	}

	if fs.props.Showing > 0 {
		selection.SetShowing(fs.props.Showing)
	}

	if err := selection.RenderWithTheme(fs.form.theme); err != nil {
		return err
	}

	fs.value = selection.Selected()

	if fs.form.scanner != nil {
		if err := fs.form.scanner(fs.value, fs.dest); err != nil {
			return err
		}
	}

	fmt.Printf("%-*s: %v\n", fs.form.maxLengthLabel, fs.label, selection.SelectedOption())

	fs.form.shownLines += 1

	return nil
}

func (fs *FormSelect) AddValidator(f func(value any) error) FormElementer {

	fs.validators = append(fs.validators, f)

	return fs
}
