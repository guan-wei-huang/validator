package validate

import (
	"fmt"
	"testing"
)

func TestNumberCompare(t *testing.T) {
	type TestData struct {
		Greater int `validate:"gt=10"`
		Less    int `validate:"lt=10"`
		Equal   int `validate:"eq=10"`
	}
	type Test struct {
		Name   string
		Datas  TestData
		Expect error
	}

	errString := "validation error: field[%s] wrong"
	cases := []Test{
		{
			Name: "positive case",
			Datas: TestData{
				Greater: 11,
				Less:    9,
				Equal:   10,
			},
			Expect: nil,
		},
		{
			Name: "great case",
			Datas: TestData{
				Greater: 4,
				Less:    9,
				Equal:   10,
			},
			Expect: fmt.Errorf(errString, "Greater"),
		},
		{
			Name: "less case",
			Datas: TestData{
				Greater: 11,
				Less:    10,
				Equal:   10,
			},
			Expect: fmt.Errorf(errString, "Less"),
		},
		{
			Name: "equal case",
			Datas: TestData{
				Greater: 11,
				Less:    10,
				Equal:   9,
			},
			Expect: fmt.Errorf(errString, "Equeal"),
		},
	}

	v := NewValidator()
	for _, tt := range cases {
		err := v.ValidateStruct(tt.Datas)
		if err != tt.Expect {
			t.Errorf("test case [%s] err, expect: %s, received: %s", tt.Name, tt.Expect, err)
		}
	}
}
