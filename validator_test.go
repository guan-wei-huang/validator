package validate

import (
	"testing"

	v10 "github.com/go-playground/validator/v10"
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

func BenchmarkNumberCompare(b *testing.B) {
	type TestData struct {
		Greater int `validate:"gt=10"`
	}
	v := NewValidator()
	for i := 0; i < b.N; i++ {
		v.ValidateStruct(TestData{
			Greater: 11,
		})
		v.ValidateStruct(TestData{
			Greater: 9,
		})
	}
}

func BenchmarkValidateNumberCompare(b *testing.B) {
	type TestData struct {
		Greater int `validate:"gt=10"`
	}
	v := v10.New()
	for i := 0; i < b.N; i++ {
		v.Struct(TestData{
			Greater: 11,
		})
		v.Struct(TestData{
			Greater: 9,
		})
	}
}
