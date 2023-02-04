package validate

import (
	"reflect"
	"strings"
)

func (v *Validator) RegisterStruct(s interface{}) error {
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
