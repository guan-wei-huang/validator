package validate

import (
	"reflect"
)

const TagName = "validate"

type Validator struct {
	ruleCache map[string]*structRule
}

func NewValidator() *Validator {
	return &Validator{}
}

// parsed validation rules for each struct
type structRule struct {
	structName   string
	structType   reflect.Type
	numFields    int
	fields       []reflect.StructField
	validateFunc [][]*validateFn
}

func newStructRule(name string, sType reflect.Type) *structRule {
	numField := sType.NumField()
	fields := make([]reflect.StructField, 0, numField)
	for i := 0; i < numField; i++ {
		fields = append(fields, sType.Field(i))
	}

	return &structRule{
		structName:   name,
		structType:   sType,
		numFields:    numField,
		fields:       fields,
		validateFunc: make([][]*validateFn, numField),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {

	return nil
}

func (v *Validator) traverseFields(vStruct interface{}, rule structRule) error {
	value := reflect.ValueOf(vStruct)

	errors := make([]error, 0)
	for i := 0; i < len(rule.fields); i++ {
		field := value.Field(i)
		for field.Kind() == reflect.Pointer && !field.IsNil() {
			field = field.Elem()
		}

		fieldValue := field.Interface()
		for _, vf := range rule.validateFunc[i] {
			if ok := vf.CheckPass(field.Kind(), fieldValue); !ok {
				// add err
				err := ErrorValidateFalse()
				errors = append(errors, err)
			}
		}
	}

	return nil
}
