package validate

import "reflect"

func interfaceToInt64(vType reflect.Kind, value interface{}) int64 {
	switch vType {
	case reflect.Int:
		return int64(value.(int))
	case reflect.Int8:
		return int64(value.(int8))
	case reflect.Int16:
		return int64(value.(int16))
	case reflect.Int32:
		return int64(value.(int32))
	}
	return 0
}
