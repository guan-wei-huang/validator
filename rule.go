package validate

import "reflect"

type applyRuleFn func(vType reflect.Kind, value, param interface{}) bool

type validateFn struct {
	fn    applyRuleFn
	param interface{}
}

func (r validateFn) CheckPass(vType reflect.Kind, v interface{}) bool {
	return r.fn(vType, v, r.param)
}

var fnTable = map[string]applyRuleFn{
	"gt": isGreater,
}

func castApplyRuleFn(funcName string, param interface{}) *validateFn {
	fn, ok := fnTable[funcName]
	if !ok {
		return nil
	}
	return &validateFn{fn, param}
}

func isGreater(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return interfaceToInt64(vType, value) > interfaceToInt64(vType, param)
	}

	return true
}
