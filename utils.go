package validate

import (
	"reflect"
	"strconv"
)

func isInt(kind reflect.Kind) bool {
	if kind >= reflect.Int && kind <= reflect.Int64 {
		return true
	}
	return false
}

func isUint(kind reflect.Kind) bool {
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		return true
	}
	return false
}

func isFloat(kind reflect.Kind) bool {
	if kind >= reflect.Float32 && kind <= reflect.Float64 {
		return true
	}
	return false
}

func parseStringToType(pType reflect.Kind, str string) (interface{}, error) {
	switch {
	case isInt(pType):
		return strconv.ParseInt(str, 10, 64)
	case isUint(pType):
		return strconv.ParseUint(str, 10, 64)
	case isFloat(pType):
		return strconv.ParseFloat(str, 64)
	}
	return nil, ErrorValidateInvalidTag()
}

func parseToInt64(vType reflect.Kind, value interface{}) int64 {
	switch vType {
	case reflect.Int:
		return int64(value.(int))
	case reflect.Int8:
		return int64(value.(int8))
	case reflect.Int16:
		return int64(value.(int16))
	case reflect.Int32:
		return int64(value.(int32))
	case reflect.Int64:
		return int64(value.(int64))
	}
	return 0
}

func paramToInt64(param interface{}) int64 {
	return param.(int64)
}

func deReferenceInterface(v interface{}) reflect.Value {
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer && !value.IsNil() {
		value = value.Elem()
	}
	return value
}
