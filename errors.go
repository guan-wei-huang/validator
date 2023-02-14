package validate

import (
	"errors"
	"fmt"
)

var (
	ErrValidatorInvalidPayload = errors.New("invalid payload")
	ErrValidatorWrongType      = errors.New("wrong type")
)

func ErrorValidateInvalidTag() error {
	return nil
}

func ErrorValidateWrongType(expect string) error {
	return fmt.Errorf("invalid validation error, expect value type: %v, but got nil", expect)
}

type ValidateError struct {
	Field string
	Rule  string
}

func (e ValidateError) Error() string {
	return fmt.Sprintf("validation failed, field: %v, violate rule: %v", e.Field, e.Rule)
}

type ValidateErrors []ValidateError

func (e ValidateErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	msg := "validation failed:\n"
	for i, v := range e {
		msg += fmt.Sprintf("%v: field[%v], violate rule[%v]\n", i, v.Field, v.Rule)
	}
	return msg
}

func ErrorValidateFalse(field, rule string) ValidateError {
	return ValidateError{field, rule}
}
