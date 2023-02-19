package validate

import (
	"reflect"
	"unsafe"
)

const TAG_NAME = "validate"

type Validator struct {
	ruleCache map[string]*structRule
}

func New() *Validator {
	return &Validator{
		ruleCache: make(map[string]*structRule),
	}
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

func (v *Validator) ValidateStruct(s interface{}) error {
	value := deReference(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}

	// register struct validate rule if cannot find rule in cache
	valueType := value.Type().String()
	if _, ok := v.ruleCache[valueType]; !ok {
		if err := v.registerStruct(value); err != nil {
			return err
		}
	}

	rule := v.ruleCache[valueType]
	return v.traverseFields(value, rule)
}

func (v *Validator) traverseFields(value reflect.Value, rule *structRule) error {
	var errors ValidateErrors

	if rule.hasUnexported {
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
		fieldValue := field.Interface()

		for _, vf := range rule.validateFunc[i] {
			if !vf.CheckPass(field.Kind(), fieldValue) {
				errors = append(errors, ErrorValidateFalse(fieldType.Name, vf.tag))
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
