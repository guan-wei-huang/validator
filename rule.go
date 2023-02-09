package validate

import "reflect"

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
	"gt": isGreater,
	"eq": isEqual,
}

func castApplyRuleFn(funcName string, param interface{}, tag string) *validateFn {
	fn, ok := fnTable[funcName]
	if !ok {
		return nil
	}
	return &validateFn{fn, param, tag}
}

func isGreater(vType reflect.Kind, value, param interface{}) bool {
	switch vType {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return parseToInt64(vType, value) > paramToInt64(param)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return
	}

	return true
}

func isEqual(vType reflect.Kind, value, param interface{}) bool {
}
