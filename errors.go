package validate

import "errors"

var (
	ErrValidatorInvalidPayload = errors.New("invalid payload")
	ErrValidatorWrongType      = errors.New("wrong type")
)
