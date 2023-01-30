package validate

import (
	"reflect"
	"strings"
)

const TagName = "validate"

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

// parsed validation rules for each struct
type structRule struct {
	structType   reflect.Type
	numFields    int
	fileds       []reflect.StructField
	validateFunc [][]*validateFn
}

// parse user defined validation rules
func parseTag(rule string) []string {
	return strings.Split(rule, ",")
}

func (v *Validator) ValidateStruct(s interface{}) error {
	val := reflect.ValueOf(s)

	for val.Kind() == reflect.Pointer && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrValidatorWrongType
	}

	valType := val.Type()
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		rule, exist := field.Tag.Lookup(TagName)
		if !exist {
			continue
		}
		rules := strings.Split(rule, ",")
		fieldName := field.Name
	}
}

func (v *Validator) traverseFields(vStruct interface{}, rule structRule) error {
	value := reflect.ValueOf(vStruct)

	errors := make([]error, 0)
	for i := 0; i < len(rule.fileds); i++ {
		field := value.Field(i)
		for field.Kind() == reflect.Pointer && !field.IsNil() {
			field = field.Elem()
		}

		fieldValue := field.Interface()
		for _, vf := range rule.validateFunc[i] {
			if ok := vf.CheckPass(field.Kind(), fieldValue); !ok {
				// add err
				var newErr error
				errors = append(errors, newErr)
			}
		}
	}
}
