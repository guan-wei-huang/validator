package validate

import (
	"errors"
	"fmt"
)

var (
	ErrValidatorInvalidPayload = errors.New("invalid payload")
	ErrValidatorWrongType      = errors.New("wrong type")
)

type ValidateError struct {
	Field string
	Rule  string
}

func (e *ValidateError) Error() string {
	return fmt.Sprintf("validation failed, field: %v, violate rule: %v", e.Field, e.Rule)
}

func ErrorValidateFalse(field, rule string) *ValidateError {
	return &ValidateError{field, rule}
}
