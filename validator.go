package validator

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

const TAG_NAME = "validate"

type Validator struct {
	ruleCache sync.Map
	// ruleCache map[string]*structRule
}

func New() *Validator {
	return &Validator{}
}

// parsed validation rules for each struct
type structRule struct {
	structName string
	structType reflect.Type

	hasUnexported bool
	fields        []reflect.StructField
	validateFunc  [][]*validateFn
}

func newStructRule(name string, sType reflect.Type) *structRule {
	hasUnexported := false
	numField := sType.NumField()
	fields := make([]reflect.StructField, numField)
	for i := 0; i < numField; i++ {
		fields[i] = sType.Field(i)
		if !fields[i].IsExported() {
			hasUnexported = true
		}
	}

	return &structRule{
		structName:    name,
		structType:    sType,
		hasUnexported: hasUnexported,
		fields:        fields,
		validateFunc:  make([][]*validateFn, numField),
	}
}

func (v *Validator) loadRule(name string) *structRule {
	if rule, ok := v.ruleCache.Load(name); ok {
		return rule.(*structRule)
	}
	return nil
}

func (v *Validator) storeRule(name string, rule *structRule) {
	v.ruleCache.Store(name, rule)
}

func (v *Validator) ValidateStruct(s interface{}) error {
	value := deref(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}

	// register struct validate rule if cannot find rule in cache
	valueType := value.Type()
	if rule := v.loadRule(valueType.String()); rule == nil || (rule != nil && rule.structType != valueType) {
		if err := v.registerStruct(valueType, valueType.String()); err != nil {
			return err
		}
	}

	rule := v.loadRule(valueType.String())
	if err := v.traverseFields(value, rule, valueType.Name()); err != nil {
		return err
	}
	return nil
}

func (v *Validator) traverseFields(value reflect.Value, rule *structRule, levelName string) ValidateErrors {
	var errors ValidateErrors

	if rule.hasUnexported {
		// reallocate an opened value
		tmp := reflect.New(value.Type()).Elem()
		tmp.Set(value)
		value = tmp
	}

	for i, fieldType := range rule.fields {
		field := value.Field(i)
		if !fieldType.IsExported() {
			field = reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr()))
		}

		for field.Kind() == reflect.Pointer && !field.IsNil() {
			field = field.Elem()
		}
		fieldKind := field.Kind()
		fieldValue := field.Interface()

		for _, vf := range rule.validateFunc[i] {
			if !vf.CheckPass(fieldKind, fieldValue) {
				name := fmt.Sprintf("%v.%v", levelName, fieldType.Name)
				errors = append(errors, ErrorValidateFalse(name, vf.tag))
			}
		}

		if fieldKind == reflect.Struct {
			nestedRule := v.loadRule(getNestedName(fieldType.Type, rule.structName, i))
			errors = append(errors, v.traverseFields(field, nestedRule, fmt.Sprintf("%v.%v", levelName, fieldType.Name))...)
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
