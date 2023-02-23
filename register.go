package validate

import (
	"reflect"
	"strings"
)

// parseTag parse tag and return slice of validateFn
func parseTag(fieldType reflect.Type, tag string, isPtr bool) ([]*validateFn, error) {
	if tag == "" {
		return nil, nil
	}

	rules := strings.Split(tag, ",")
	fs := make([]*validateFn, 0, len(rules))
	for _, r := range rules {
		name, param, _ := strings.Cut(r, "=")
		switch name {
		case "gt", "eq", "ls":
			p, err := parseStringToType(fieldType.Kind(), param)
			if err != nil {
				return nil, err
			}
			fs = append(fs, castApplyRuleFn(name, p, r))
		case "len":
			p, err := parseStringToType(reflect.Int, param)
			if err != nil {
				return nil, err
			}
			fs = append(fs, castApplyRuleFn(name, p, r))
		case "required":
			fs = append(fs, castApplyRuleFn(name, isPtr, r))
		default:
		}
	}

	return fs, nil
}

func (v *Validator) RegisterStruct(s interface{}) error {
	value := deReference(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}
	return v.registerStruct(value.Type())
}

func (v *Validator) registerStruct(vType reflect.Type, name ...string) error {
	ruleName := vType.String()
	if len(name) != 0 {
		ruleName = name[0]
	}
	rule := newStructRule(ruleName, vType)

	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		fieldType := field.Type

		// dereference if field is pointer
		isPtr := false
		for fieldType.Kind() == reflect.Pointer {
			isPtr = true
			fieldType = fieldType.Elem()
		}

		tag, _ := field.Tag.Lookup(TAG_NAME)
		fs, err := parseTag(fieldType, tag, isPtr)
		if err != nil {
			return err
		}
		rule.validateFunc[i] = fs

		// register for nested struct
		if fieldType.Kind() == reflect.Struct {
			nestedName := getNestedName(fieldType, ruleName, i)
			if err := v.registerStruct(fieldType, nestedName); err != nil {
				return err
			}
		}
	}

	// push into cache
	v.storeRule(ruleName, rule)
	return nil
}
