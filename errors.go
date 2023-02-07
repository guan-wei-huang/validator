package validate

import (
	"errors"
	"fmt"
)

var (
	ErrValidatorInvalidPayload = errors.New("invalid payload")
	ErrValidatorWrongType      = errors.New("wrong type")
)

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

func ErrorValidateFalse(field, rule string) ValidateError {
	return ValidateError{field, rule}
}
