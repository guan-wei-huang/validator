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
		name, param, _ := strings.Cut(r, ":")
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
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Pointer && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrValidatorWrongType
	}

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
	v.ruleCache[valType.String()] = rule
	return nil
}
