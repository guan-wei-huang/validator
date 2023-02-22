package validate

import (
	"reflect"
)

type applyRuleFn func(vType reflect.Kind, value, param interface{}) bool

type validateFn struct {
	fn    applyRuleFn
	param interface{}
	tag   string
}

func (r validateFn) CheckPass(vType reflect.Kind, v interface{}) bool {
	return r.fn(vType, v, r.param)
}

var fnTable = map[string]applyRuleFn{
	"gt":       isGreater,
	"eq":       isEqual,
	"ls":       isLess,
	"len":      isLen,
	"required": isRequired,
}

func castApplyRuleFn(funcName string, param interface{}, tag string) *validateFn {
	fn, ok := fnTable[funcName]
	if !ok {
		return nil
	}
	return &validateFn{fn, param, tag}
}

// if vType is excluded in switch case, it must be reflect.Pointer.
// happen when the field is pointer type and user given value is a nil pointer.
func isGreater(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseToInt64(vType, value) > param.(int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseToUint64(vType, value) > param.(uint64)
	case reflect.Float32, reflect.Float64:
		return parseToFloat64(vType, value) > param.(float64)
	}
	return false
}

func isEqual(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseToInt64(vType, value) == param.(int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseToUint64(vType, value) == param.(uint64)
	case reflect.Float32, reflect.Float64:
		return parseToFloat64(vType, value) == param.(float64)
	}
	return false
}

func isLess(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseToInt64(vType, value) < param.(int64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return parseToUint64(vType, value) < param.(uint64)
	case reflect.Float32, reflect.Float64:
		return parseToFloat64(vType, value) < param.(float64)
	}
	return false
}

func isLen(vType reflect.Kind, value, param interface{}) bool {
	size := int(param.(int64))
	switch vType {
	case reflect.String:
		return len(value.(string)) == size
	case reflect.Array:
		return reflect.ValueOf(value).Len() == size
	case reflect.Slice:
		return reflect.ValueOf(value).Len() == size
	}
	return false
}

// isRequired param stores whether origin field's type is pointer or not.
// if it's ptr, verify that value is not nil. otherwise, check that value is not
// empty value
func isRequired(vType reflect.Kind, value, param interface{}) bool {
	if isPtr := param.(bool); isPtr {
		return vType != reflect.Pointer
	}

	return !reflect.ValueOf(value).IsZero()
}
