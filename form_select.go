/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"fmt"
	"strings"
)

type SelectProps struct {
	Options          []string
	Values           []any
	Showing          int
	Selected         func(value any) bool
	InfoText         string
	OptionsAndValues func() ([]string, []any, error)
	Callback         func(value any) string
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

	if fs.props.OptionsAndValues != nil {
		var err error
		fs.props.Options, fs.props.Values, err = fs.props.OptionsAndValues()
		if err != nil {
			return err
		}
	}

	var clearLines int
	if fs.props.InfoText != "" {
		fmt.Println(fs.props.InfoText)
		clearLines = 1 + strings.Count(fs.props.InfoText, "\n")
	}

	selection, err := NewSelection(fs.props.Options, fs.props.Values)
	if err != nil {
		return err
	}

	if fs.defaultValue != nil {
		for p, v := range fs.props.Values {
			if fs.defaultValue(&DefaultValueProps{SelectOption: v}).Found {
				selection.SetSelected(p)
				break
			}
		}
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

	fs.form.shownLines += clearLines

	return nil
}

func (fs *FormSelect) AddValidator(f func(value any) error) FormElementer {

	fs.validators = append(fs.validators, f)

	return fs
}

func (fs *FormSelect) DefaultValue(f func(props *DefaultValueProps) DefaultValue) FormElementer {

	fs.defaultValue = f

	return fs
}

func (fs *FormSelect) Callback() {
	if fs.props.Callback != nil {
		if fs.props.InfoText != "" {
			ClearLines(2 + strings.Count(fs.props.InfoText, "\n"))
		}

		fmt.Println(fs.props.Callback(fs.value))
	}
}
