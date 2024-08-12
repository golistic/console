/*
 * Copyright (c) 2024, Geert JM Vanderkelen
 */

package console

func NewForm() *Form {
	return &Form{}
}

func NewFormWithScanner(scanner func(value any, dest any) error) *Form {
	return &Form{
		scanner: scanner,
		theme:   "ascii",
	}
}

type Form struct {
	Elements []FormElementer
	scanner  func(value any, dest any) error

	maxLengthLabel int
	theme          Theme
	shownLines     int
}

func (f *Form) SetTheme(theme Theme) *Form {
	f.theme = theme
	return f
}

func (f *Form) AddElements(elements ...FormElementer) {

	f.Elements = []FormElementer{}

	for _, element := range elements {
		element.setForm(f)
		f.Elements = append(f.Elements, element)
	}
}

func (f *Form) Execute() error {

	for _, elm := range f.Elements {
		l := len(elm.Label())
		if l > f.maxLengthLabel {
			f.maxLengthLabel = l
		}
	}

	for _, element := range f.Elements {
		if err := element.do(); err != nil {
			return err
		}

		for _, validator := range element.getValidators() {
			if err := validator(element.Value()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Form) Clear() {
	ClearLines(f.shownLines + 1)
}

func (f *Form) RawValues() map[string]any {

	values := map[string]any{}

	for _, elm := range f.Elements {
		values[elm.Name()] = elm.Value()
	}

	return values
}
