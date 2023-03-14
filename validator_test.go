package validator

import (
	"reflect"
	"testing"

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

func combineValidateError(fields, rules []string) string {
	if len(fields) != len(rules) {
		return ""
	}
	es := make(ValidateErrors, 0, len(rules))
	for i := range fields {
		err := ErrorValidateFalse(fields[i], rules[i])
		es = append(es, err)
	}
	return es.Error()
}

func TestValidateWrongType(t *testing.T) {
	validate := New()
	err := validate.ValidateStruct("wrong type")
	assert.EqualError(t, err, ErrorValidateWrongType(reflect.Struct.String()).Error())
}

func TestNumberCompare(t *testing.T) {
	type TestData struct {
		intGt     int       `validate:"gt=10"`
		IntPtrGt  *int      `validate:"gt=10"`
		FloatGt   float32   `validate:"gt=10.1"`
		IntEq     int       `validate:"eq=10"`
		IntLs     int       `validate:"ls=10"`
		Int8Eq    int8      `validate:"eq=10"`
		Int16Eq   int16     `validate:"eq=10"`
		Int32Eq   int32     `validate:"eq=10"`
		Int64Eq   int64     `validate:"eq=10"`
		UintEq    uint      `validate:"eq=1"`
		Uint64Eq  uint64    `validate:"eq=2"`
		ComplexEq complex64 `validate:"eq=10+2i"`
	}

	validate := New()
	err := validate.ValidateStruct(TestData{
		intGt:     11,
		IntPtrGt:  toPtr(11),
		FloatGt:   10.3,
		IntEq:     10,
		IntLs:     9,
		Int8Eq:    10,
		Int16Eq:   10,
		Int32Eq:   10,
		Int64Eq:   10,
		UintEq:    1,
		Uint64Eq:  2,
		ComplexEq: 10 + 2i,
	})
	assert.ErrorIs(t, err, nil)

	err = validate.ValidateStruct(TestData{
		intGt:     9,
		IntPtrGt:  toPtr(9),
		FloatGt:   9.2,
		IntEq:     9,
		IntLs:     11,
		Int8Eq:    6,
		Int16Eq:   7,
		Int32Eq:   8,
		Int64Eq:   9,
		UintEq:    2,
		Uint64Eq:  3,
		ComplexEq: 10,
	})
	assert.EqualError(t, err, combineValidateError(
		[]string{"TestData.intGt", "TestData.IntPtrGt", "TestData.FloatGt", "TestData.IntEq", "TestData.IntLs", "TestData.Int8Eq", "TestData.Int16Eq",
			"TestData.Int32Eq", "TestData.Int64Eq", "TestData.UintEq", "TestData.Uint64Eq", "TestData.ComplexEq"},
		[]string{"gt=10", "gt=10", "gt=10.1", "eq=10", "ls=10", "eq=10", "eq=10", "eq=10", "eq=10", "eq=1", "eq=2", "eq=10+2i"},
	))
}

func TestLen(t *testing.T) {
	type TestData struct {
		Str        string    `validate:"len=4"`
		FloatSlice []float32 `validate:"len=2"`
		UintArray  [2]uint   `validate:"len=2"`
	}

	validate := New()
	err := validate.ValidateStruct(TestData{
		Str:        "test",
		FloatSlice: []float32{2.4, 3.7},
		UintArray:  [2]uint{3, 5},
	})
	assert.NoError(t, err)

	err = validate.ValidateStruct(TestData{
		Str:        "wrong",
		FloatSlice: []float32{1.3},
		UintArray:  [2]uint{7},
	})
	assert.EqualError(t, err, combineValidateError(
		[]string{"TestData.Str", "TestData.FloatSlice"},
		[]string{"len=4", "len=2"},
	))
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
	assert.EqualError(t, err, combineValidateError(
		[]string{"TestData.Num", "TestData.Str", "TestData.NumSlice", "TestData.NumPtr", "TestData.StrPtr",
			"TestData.M", "TestData.C"},
		[]string{"required", "required", "required", "required", "required", "required", "required"},
	))
}

func TestMinMax(t *testing.T) {
	t.Skip()
	type TestData struct {
		MinSlice []int     `validate:"min=2"`
		MaxSlice []float32 `validate:"max=100.2"`
		Array    [2]int    `validate:"min=10,max=30"`
	}

	validate := New()
	err := validate.ValidateStruct(TestData{
		MinSlice: []int{3, 4, 6},
		MaxSlice: []float32{10.2, 40.2},
		Array:    [2]int{13, 30},
	})
	assert.NoError(t, err)

	err = validate.ValidateStruct(TestData{
		MinSlice: []int{5, 6, 1},
		MaxSlice: []float32{20, 30, 1000},
		Array:    [2]int{3, 41},
	})
	assert.EqualError(t, err, combineValidateError(
		[]string{"MinSlice", "MaxSlice", "Array"},
		[]string{"min=2", "max=100.2", "min=10"},
	))
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

func TestRepeatNameStruct(t *testing.T) {
	validate := New()
	type TestData struct {
		Num int `validate:"gt=10,required"`
	}
	err := validate.ValidateStruct(TestData{
		Num: 3,
	})
	assert.EqualError(t, err, combineValidateError([]string{"TestData.Num"}, []string{"gt=10"}))

	t.Run("different struct with the same name", func(t *testing.T) {
		type TestData struct {
			Str string `validate:"len=3,required"`
		}
		err := validate.ValidateStruct(TestData{
			Str: "test",
		})
		assert.EqualError(t, err, combineValidateError([]string{"TestData.Str"}, []string{"len=3"}))
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
		assert.EqualError(t, err, combineValidateError(
			[]string{"TestData.Num", "TestData.Nested.Num", "TestData.Nested.Str"},
			[]string{"gt=10", "gt=10", "len=4"},
		))
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
				Num int    `validate:"gt=10"`
				Str string `validate:"required"`
			}{
				Num: 2,
				Str: "",
			},
		})
		assert.EqualError(t, err, combineValidateError(
			[]string{"UndefNested.Nested.Num", "UndefNested.Nested.Str"},
			[]string{"gt=10", "required"},
		))
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
		assert.EqualError(t, err, combineValidateError(
			[]string{"Case.Num", "Case.Nested.str", "Case.Nested.intPtr", "Case.Nested.intSlice"},
			[]string{"gt=10", "len=4", "eq=10", "len=2"},
		))
	})
}

func TestRegister(t *testing.T) {
	type TestData struct {
		Num int    `validate:"eq=2"`
		Str string `validate:"len=3"`
	}
	validate := New()

	err := validate.RegisterStruct("wrong type")
	assert.EqualError(t, err, ErrorValidateWrongType(reflect.Struct.String()).Error())

	err = validate.RegisterStruct(&TestData{})
	assert.NoError(t, err)

	err = validate.ValidateStruct(&TestData{
		Num: 1,
		Str: "test",
	})
	assert.EqualError(t, err, combineValidateError(
		[]string{"TestData.Num", "TestData.Str"},
		[]string{"eq=2", "len=3"},
	))
}

func TestRegisterByMap(t *testing.T) {
	type Nested struct {
		Nint int
		Nstr string
	}
	type TestData struct {
		Num      int
		Str      string
		IntSlice []int
		M        map[int]int
		Nested   Nested
	}

	validate := New()

	t.Run("success", func(t *testing.T) {
		err := validate.RegisterMapRule(TestData{}, map[string]interface{}{
			"Num":      "gt=10,required",
			"Str":      "len=4,required",
			"IntSlice": "len=2",
			"M":        "required",
			"Nested": map[string]interface{}{
				"Nint": "eq=3",
				"Nstr": "len=2,required",
			},
		})
		assert.NoError(t, err)

		err = validate.ValidateStruct(TestData{
			Num:      11,
			Str:      "test",
			IntSlice: []int{3, 4},
			M:        make(map[int]int),
			Nested: Nested{
				Nint: 3,
				Nstr: "tt",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("failed", func(t *testing.T) {
		err := validate.RegisterMapRule(TestData{}, map[string]interface{}{
			"Num":      "gt=10,required",
			"Str":      "len=4,required",
			"IntSlice": "len=2",
			"M":        "required",
			"Nested": map[string]interface{}{
				"Nint": "eq=3",
				"Nstr": "len=2,required",
			},
		})
		assert.NoError(t, err)

		err = validate.ValidateStruct(TestData{
			Num:      8,
			Str:      "test",
			IntSlice: []int{3, 4},
			M:        make(map[int]int),
			Nested: Nested{
				Nint: 2,
				Nstr: "test",
			},
		})
		assert.EqualError(t, err, combineValidateError(
			[]string{"TestData.Num", "TestData.Nested.Nint", "TestData.Nested.Nstr"},
			[]string{"gt=10", "eq=3", "len=2"},
		))
	})

	t.Run("register failed", func(t *testing.T) {
		err := validate.RegisterMapRule(TestData{}, map[string]interface{}{
			"Num": "gt=10,unsupported",
			"Str": "len=4",
			"M":   "required",
		})
		assert.EqualError(t, err, ErrorValidateUnsupportedTag("unsupported").Error())
	})
}
