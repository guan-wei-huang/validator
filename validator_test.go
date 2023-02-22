package validate

import (
	"testing"

	v10 "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestUnexported struct {
	num      int         `validate:"required,gt=10"`
	str      string      `validate:"required,len=4"`
	intPtr   *int        `validate:"required,eq=10"`
	strPtr   *string     `validate:"required"`
	intSlice []int       `validate:"required,len=2"`
	m        map[int]int `validate:"required"`
	c        chan int    `validate:"required"`
}

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

func TestNumberCompare(t *testing.T) {
	type TestData struct {
		intGt    int     `validate:"gt=10"`
		IntPtrGt *int    `validate:"gt=10"`
		FloatGt  float32 `validate:"gt=10.1"`
		IntEq    int     `validate:"eq=10"`
		IntLs    int     `validate:"ls=10"`
	}

	validate := New()
	err := validate.ValidateStruct(TestData{
		intGt:    11,
		IntPtrGt: toPtr(11),
		FloatGt:  10.3,
		IntEq:    10,
		IntLs:    9,
	})
	assert.ErrorIs(t, err, nil)

	err = validate.ValidateStruct(TestData{
		intGt:    9,
		IntPtrGt: toPtr(9),
		FloatGt:  9.2,
		IntEq:    9,
		IntLs:    11,
	})
	assert.EqualError(t, err, combinateValidateError(
		[]string{"intGt", "IntPtrGt", "FloatGt", "IntEq", "IntLs"},
		[]string{"gt=10", "gt=10", "gt=10.1", "eq=10", "ls=10"},
	).Error())
}

func TestRequired(t *testing.T) {
	type TestData struct {
		Num      int         `validate:"required"`
		Str      string      `validate:"required"`
		NumSlice []int       `validate:"required"`
		NumPtr   *int        `validate:"required"`
		StrPtr   *string     `validate:"required"`
		M        map[int]int `validate:"required"`
		C        chan int    `validate:"required"`
	}

	validate := New()
	err := validate.ValidateStruct(TestData{
		Num:      1,
		Str:      "test",
		NumSlice: []int{2, 3},
		NumPtr:   toPtr(0),
		StrPtr:   toPtr(""),
		M:        map[int]int{1: 2},
		C:        make(chan int),
	})
	assert.ErrorIs(t, err, nil)

	err = validate.ValidateStruct(TestData{
		Num: 0,
		Str: "",
	})
	assert.EqualError(t, err, combinateValidateError(
		[]string{"Num", "Str", "NumSlice", "NumPtr", "StrPtr", "M", "C"},
		[]string{"required", "required", "required", "required", "required", "required", "required"},
	).Error())
}

func TestUnexportedField(t *testing.T) {
	validate := New()
	assert.NotPanics(t, func() {
		validate.ValidateStruct(TestUnexported{
			num:      11,
			str:      "xxx",
			intPtr:   toPtr(2),
			strPtr:   toPtr("xx"),
			intSlice: []int{1},
			m:        make(map[int]int),
			c:        make(chan int, 1),
		})
	})
}

func TestNestedStruct(t *testing.T) {
	type NestedStruct struct {
		Num int    `validate:"gt=10,required"`
		Str string `validate:"len=4,required"`
	}
	type TestData struct {
		Num    int          `validate:"gt=10"`
		Nested NestedStruct `validate:"required"`
	}

	validate := New()
	t.Run("success case", func(t *testing.T) {
		err := validate.ValidateStruct(TestData{
			Num: 11,
			Nested: NestedStruct{
				Num: 11,
				Str: "test",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("failed inside nested struct", func(t *testing.T) {
		err := validate.ValidateStruct(TestData{
			Num: 9,
			Nested: NestedStruct{
				Num: 9,
				Str: "small",
			},
		})
		assert.EqualError(t, err, combinateValidateError(
			[]string{"validate.TestData.Num", "validate.NestedStruct.Num", "validate.NestedStruct.Str"},
			[]string{"gt=10", "gt=10", "len=4"},
		).Error())
	})

	t.Run("undefined nested struct", func(t *testing.T) {
		type UndefNested struct {
			Num    int `validate:"gt=10"`
			Nested struct {
				Num int    `validate:"gt=10"`
				Str string `validate:"required"`
			}
		}
		err := validate.ValidateStruct(UndefNested{
			Num: 11,
			Nested: struct {
				Num int    "validate:\"gt=10\""
				Str string "validate:\"required\""
			}{
				Num: 2,
				Str: "",
			},
		})
		assert.EqualError(t, err, combinateValidateError(
			[]string{"validate.UndefNested-1.Num", "validate.UndefNested-1.Str"},
			[]string{"gt=10", "required"},
		).Error())
	})

	t.Run("unexported field in nested struct", func(t *testing.T) {
		type Case struct {
			Num    int            `validate:"gt=10"`
			Nested TestUnexported `validate:"required"`
		}
		err := validate.ValidateStruct(Case{
			Num: 8,
			Nested: TestUnexported{
				num:      11,
				str:      "xxx",
				intPtr:   toPtr(2),
				strPtr:   toPtr("xx"),
				intSlice: []int{1},
				m:        make(map[int]int),
				c:        make(chan int, 1),
			},
		})
		assert.EqualError(t, err, combinateValidateError(
			[]string{"validate.Case.Num", "validate.TestUnexported.str", "validate.TestUnexported.intPtr",
				"validate.TestUnexported.intSlice"},
			[]string{"gt=10", "len=4", "eq=10", "len=2"},
		).Error())
	})
}

func BenchmarkPackage(b *testing.B) {
	type TestData struct {
		IntGt   int  `validate:"gt=10"`
		IntEq   int  `validate:"eq=10"`
		Require *int `validate:"required"`
	}
	validate := New()
	for i := 0; i < b.N; i++ {
		validate.ValidateStruct(TestData{
			IntGt:   12,
			IntEq:   10,
			Require: toPtr(2),
		})
	}
}

func BenchmarkV10(b *testing.B) {
	type TestData struct {
		IntGt   int  `validate:"gt=10"`
		IntEq   int  `validate:"eq=10"`
		Require *int `validate:"required"`
	}
	v := v10.New()
	for i := 0; i < b.N; i++ {
		v.Struct(TestData{
			IntGt:   12,
			IntEq:   10,
			Require: toPtr(2),
		})
	}
}
