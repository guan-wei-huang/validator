package validate

import (
	"reflect"
	"strings"
)

// parseTag parse tag and return slice of validateFn
func parseTag(fieldType reflect.Type, tag string, isPtr bool) ([]*validateFn, error) {
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
		case "required":
			fs = append(fs, castApplyRuleFn(name, isPtr, r))
		default:
		}
	}

	return fs, nil
}

func (v *Validator) RegisterStruct(s interface{}) error {
	value := deReferenceInterface(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}
	return v.registerStruct(value)
}

func (v *Validator) registerStruct(val reflect.Value) error {
	valType := val.Type()
	rule := newStructRule(valType.String(), valType)
	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		tag, exist := field.Tag.Lookup(TAG_NAME)
		if !exist {
			continue
		}

		fieldType := field.Type
		isPtr := false
		for fieldType.Kind() == reflect.Pointer {
			isPtr = true
			fieldType = fieldType.Elem()
		}

		fs, err := parseTag(fieldType, tag, isPtr)
		if err != nil {
			return err
		}
		rule.validateFunc[i] = fs
	}

	// push into cache
	v.ruleCache[val.Type().String()] = rule
	return nil
}
