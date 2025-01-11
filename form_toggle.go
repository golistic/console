/*
 * Copyright (c) 2025, Geert JM Vanderkelen
 */

package console

type ToggleProps struct {
	Options      []string
	Values       []any
	Selected     func(value any) bool
	DefaultValue any
}

func NewFormToggle(name, label string, dest any, props ToggleProps) *FormToggle {
	return &FormToggle{
		formElement: &formElement{
			name:  name,
			label: label,
			dest:  dest,
		},

		props: props,
	}
}

func NewFormToggleBool(name, label string, dest *bool, defaultValue any) *FormToggle {
	return NewFormToggle(name, label, dest, ToggleProps{
		Options:      []string{"Yes", "No"},
		Values:       []any{true, false},
		DefaultValue: defaultValue,
	})
}

type FormToggle struct {
	*formElement

	props ToggleProps
}

var _ FormElementer = (*FormToggle)(nil)

func (ft *FormToggle) do() error {

	toggle, err := NewToggle(ft.label, ft.props.Options, ft.props.Values)
	if err != nil {
		return err
	}

	toggle.SetSelected(ft.props.DefaultValue)

	if err := toggle.Render(); err != nil {
		return err
	}

	ft.value = toggle.Selected()

	if ft.form.scanner != nil {
		return ft.form.scanner(ft.value.(bool), ft.dest)
	}

	return nil
}

func (ft *FormToggle) AddValidator(f func(value any) error) FormElementer {

	ft.validators = append(ft.validators, f)

	return ft
}

func (ft *FormToggle) DefaultValue(f func(props *DefaultValueProps) DefaultValue) FormElementer {

	ft.defaultValue = f

	return ft
}
