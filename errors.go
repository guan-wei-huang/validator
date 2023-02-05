package validate

import "errors"

var (
	ErrValidatorInvalidPayload = errors.New("invalid payload")
	ErrValidatorWrongType      = errors.New("wrong type")
)

type ValidateError struct {
	Field string
	Rule  string
}

func (e *ValidateError) Error() string {
	return ""
}

func ErrorValidateFalse() *ValidateError {
	return &ValidateError{}
}
