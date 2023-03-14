package validator

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
		var vfn *validateFn
		name, param, _ := strings.Cut(r, "=")
		switch name {
		case "gt", "eq", "ls":
			p, err := parseStringToType(fieldType.Kind(), param)
			if err != nil {
				return nil, err
			}
			vfn = castApplyRuleFn(name, p, r)

		case "min", "max":
			if !isArrayBased(fieldType.Kind()) {
				return nil, ErrorValidateUnsupportedTag(r)
			}
			ftype := fieldType.Elem()
			p, err := parseStringToType(ftype.Kind(), param)
			if err != nil {
				return nil, err
			}
			vfn = castApplyRuleFn(name, p, r)

		case "len":
			p, err := parseStringToType(reflect.Int, param)
			if err != nil {
				return nil, err
			}
			vfn = castApplyRuleFn(name, p, r)

		case "required":
			vfn = castApplyRuleFn(name, isPtr, r)

		default:
			return nil, ErrorValidateUnsupportedTag(r)
		}
		fs = append(fs, vfn)
	}

	return fs, nil
}

func (v *Validator) RegisterMapRule(s interface{}, ruleMap map[string]interface{}) error {
	value := deref(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}
	vType := value.Type()
	return v.registerMapRule(vType, ruleMap, vType.String())
}

func (v *Validator) registerMapRule(vType reflect.Type, ruleMap map[string]interface{}, ruleName string) error {
	rule := newStructRule(ruleName, vType)
	for i := 0; i < vType.NumField(); i++ {
		field := vType.Field(i)
		fieldType := field.Type
		fieldRule, exist := ruleMap[field.Name]
		if !exist {
			continue
		}

		isPtr := false
		for fieldType.Kind() == reflect.Pointer {
			isPtr = true
			fieldType = fieldType.Elem()
		}

		// nested map rule
		if nestedRule, ok := fieldRule.(map[string]interface{}); ok {
			nestedName := getNestedName(fieldType, ruleName, i)
			if err := v.registerMapRule(fieldType, nestedRule, nestedName); err != nil {
				return err
			}
		}
		if strRule, ok := fieldRule.(string); ok {
			fs, err := parseTag(fieldType, strRule, isPtr)
			if err != nil {
				return err
			}
			rule.validateFunc[i] = fs
		}
	}

	v.storeRule(ruleName, rule)
	return nil
}

func (v *Validator) RegisterStruct(s interface{}) error {
	value := deref(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}
	vType := value.Type()
	return v.registerStruct(vType, vType.String())
}

func (v *Validator) registerStruct(vType reflect.Type, ruleName string) error {
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
