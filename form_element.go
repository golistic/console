/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

type FormElementer interface {
	setForm(*Form)
	do() error
	getValidators() []func(value any) error

	Label() string
	Name() string
	Value() any
	AddValidator(func(value any) error) FormElementer
}

type formElement struct {
	name       string
	label      string
	dest       any
	value      any
	form       *Form
	validators []func(value any) error
}

var _ FormElementer = (*formElement)(nil)

func (fe *formElement) Value() any {

	return fe.value
}

func (fe *formElement) setForm(f *Form) {

	fe.form = f
}

func (fe *formElement) do() error {

	panic("cannot use formElement directly")
}

func (fe *formElement) Label() string {

	return fe.label
}

func (fe *formElement) getValidators() []func(value any) error {

	return fe.validators
}

func (fe *formElement) Name() string {

	return fe.name
}

func (fe *formElement) AddValidator(f func(value any) error) FormElementer {

	panic("need to implement AddValidator")
}
