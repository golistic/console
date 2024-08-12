/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

import (
	"fmt"

	"github.com/ergochat/readline"
)

func NewFormInput(name, label string, dest any) *FormInput {
	return &FormInput{
		formElement: &formElement{
			name:  name,
			label: label,
			dest:  dest,
		},
	}
}

type FormInput struct {
	*formElement
}

var _ FormElementer = (*FormInput)(nil)

func (fi *FormInput) do() error {

	rl, err := readline.New("")
	if err != nil {
		return err
	}
	defer func() { _ = rl.Close() }()

	rl.SetPrompt(fmt.Sprintf("%-*s: ", fi.form.maxLengthLabel, fi.label))

	fi.value, err = rl.Readline()
	if err != nil {
		return err
	}

	fi.form.shownLines += 1

	if fi.form.scanner != nil {
		return fi.form.scanner(fi.value.(string), fi.dest)
	}

	return nil
}

func (fi *FormInput) AddValidator(f func(value any) error) FormElementer {

	fi.validators = append(fi.validators, f)

	return fi
}