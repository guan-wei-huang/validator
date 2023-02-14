package validate

import (
	"testing"

	v10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func toPtr[T int | string](v T) *T {
	return &v
}

func combinateValidateError(fields, rules []string) error {
	if len(fields) != len(rules) {
		return nil
	}
	es := make(ValidateErrors, 0, len(rules))
	for i := range fields {
		err := ErrorValidateFalse(fields[i], rules[i])
		es = append(es, err)
	}
	return es
}

func TestIntValidation(t *testing.T) {
	type TestData struct {
		IntGt    int     `validate:"gt=10"`
		IntPtrGt *int    `validate:"gt=10"`
		FloatGt  float32 `validate:"gt=10.1"`
		IntEq    int     `validate:"eq=10"`
		IntLs    int     `validate:"ls=10"`
	}

	v := NewValidator()
	err := v.ValidateStruct(TestData{
		IntGt:    11,
		IntPtrGt: toPtr(11),
		FloatGt:  10.3,
		IntEq:    10,
		IntLs:    9,
	})
	assert.ErrorIs(t, err, nil)

	err = v.ValidateStruct(TestData{
		IntGt:    9,
		IntPtrGt: toPtr(9),
		FloatGt:  9.2,
		IntEq:    9,
		IntLs:    11,
	})
	assert.EqualError(t, err, combinateValidateError(
		[]string{"IntGt", "IntPtrGt", "FloatGt", "IntEq", "IntLs"},
		[]string{"gt=10", "gt=10", "gt=10.1", "eq=10", "ls=10"},
	).Error())
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
