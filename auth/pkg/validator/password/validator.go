package password

import "errors"

type Validator struct {
	ValidationFuncs ValidationFuncs
}

func NewValidator(vf ...ValidationFunc) *Validator {
	v := &Validator{ValidationFuncs: make(ValidationFuncs, 0, len(vf))}
	for _, f := range vf {
		v.ValidationFuncs = append(v.ValidationFuncs, f)
	}
	return v
}

func (v *Validator) Validate(password string) error {
	for _, f := range v.ValidationFuncs {
		if err := f(password); err != nil {
			return err
		}
	}
	return nil
}

type ValidationFunc func(password string) error

type ValidationFuncs []func(password string) error

type DefaultValidationFuncs struct {
	MinPassLength int
	MaxPassLength int
}

var (
	ErrInvalidValidationFuncsConfig = errors.New("invalid validation funcs config")

	ErrPasswordIsTooShort = errors.New("password is too short")
	ErrPasswordIsTooLong  = errors.New("password is too long")
)

func NewDefaultValidationFuncs(minPassLength, maxPassLength int) (*DefaultValidationFuncs, error) {
	if maxPassLength != 0 && maxPassLength-minPassLength <= 0 {
		return nil, ErrInvalidValidationFuncsConfig
	}

	vf := &DefaultValidationFuncs{
		MinPassLength: minPassLength,
		MaxPassLength: maxPassLength,
	}
	return vf, nil
}

func (vf *DefaultValidationFuncs) ValidatePasswordLength(password string) error {
	switch {
	case len(password) < vf.MinPassLength:
		return ErrPasswordIsTooShort
	case len(password) > vf.MaxPassLength:
		return ErrPasswordIsTooLong
	}
	return nil
}
