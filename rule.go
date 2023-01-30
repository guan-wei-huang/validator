package validate

import "reflect"

type ApplyRuleFn func(vType reflect.Kind, value, param interface{}) bool

type validateFn struct {
	fn    ApplyRuleFn
	param interface{}
}

func (r validateFn) CheckPass(vType reflect.Kind, v interface{}) bool {
	return r.fn(vType, v, r.param)
}

func isGreater(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return interfaceToInt64(vType, value) > interfaceToInt64(vType, param)
	}

	return true
}
