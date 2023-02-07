package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberCompare(t *testing.T) {
	type TestData struct {
		Greater int `validate:"gt=10"`
	}

	v := NewValidator()
	err := v.ValidateStruct(TestData{
		Greater: 11,
	})
	assert.ErrorIs(t, err, nil)

	err = v.ValidateStruct(TestData{
		Greater: 9,
	})
	assert.ErrorIs(t, err, ErrorValidateFalse("Greater", "gt=10"))
}
