package validate

import (
	"reflect"
	"strconv"
	"strings"
)

// parse tag and return slice of validateFn
func parseTag(tag string) ([]*validateFn, error) {
	rules := strings.Split(tag, ",")
	fs := make([]*validateFn, 0, len(rules))
	for _, r := range rules {
		name, param, _ := strings.Cut(r, "=")
		switch name {
		case "gt", "eq":
			paramI64, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return nil, err
			}
			fs = append(fs, castApplyRuleFn(name, paramI64, r))
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
		tag, exist := field.Tag.Lookup(TagName)
		if !exist {
			continue
		}

		fs, err := parseTag(tag)
		if err != nil {
			return err
		}
		rule.validateFunc[i] = fs
	}

	// push into cache
	v.ruleCache[val.Type().String()] = rule
	// log.Printf("%+v", rule)
	return nil
}
