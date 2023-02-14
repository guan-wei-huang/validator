package validate

import (
	"reflect"
)

const TAG_NAME = "validate"

type Validator struct {
	ruleCache map[string]*structRule
}

func NewValidator() *Validator {
	return &Validator{
		ruleCache: make(map[string]*structRule),
	}
}

// parsed validation rules for each struct
type structRule struct {
	structName   string
	structType   reflect.Type
	numFields    int
	fields       []reflect.StructField
	validateFunc [][]*validateFn
}

func newStructRule(name string, sType reflect.Type) *structRule {
	numField := sType.NumField()
	fields := make([]reflect.StructField, numField)
	for i := 0; i < numField; i++ {
		fields[i] = sType.Field(i)
	}

	return &structRule{
		structName:   name,
		structType:   sType,
		numFields:    numField,
		fields:       fields,
		validateFunc: make([][]*validateFn, numField),
	}
}

func (v *Validator) ValidateStruct(s interface{}) error {
	value := deReferenceInterface(s)
	if value.Kind() != reflect.Struct {
		return ErrorValidateWrongType(reflect.Struct.String())
	}

	// register struct validate rule if doesnt find rule in cache
	if _, ok := v.ruleCache[value.Type().String()]; !ok {
		if err := v.registerStruct(value); err != nil {
			return err
		}
	}

	rule := v.ruleCache[value.Type().String()]
	return v.traverseFields(value, rule)
}

func (v *Validator) traverseFields(value reflect.Value, rule *structRule) error {
	var errors ValidateErrors
	for i := 0; i < len(rule.fields); i++ {
		field := value.Field(i)
		for field.Kind() == reflect.Pointer && !field.IsNil() {
			field = field.Elem()
		}

		fieldValue := field.Interface()
		for _, vf := range rule.validateFunc[i] {
			if ok := vf.CheckPass(field.Kind(), fieldValue); !ok {
				errors = append(errors, ErrorValidateFalse(rule.fields[i].Name, vf.tag))
			}
		}
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}
