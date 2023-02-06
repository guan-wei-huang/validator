package validate

import (
	"reflect"
)

const TagName = "validate"

type Validator struct {
	ruleCache map[string]*structRule
}

func NewValidator() *Validator {
	return &Validator{}
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
	fields := make([]reflect.StructField, 0, numField)
	for i := 0; i < numField; i++ {
		fields = append(fields, sType.Field(i))
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
	if value.Kind() != reflect.Struct || value.IsNil() {
		return ErrValidatorWrongType
	}

	// register struct validate rule if doesnt find rule in cache
	if _, ok := v.ruleCache[value.Type().Name()]; !ok {
		if err := v.RegisterStruct(s); err != nil {
			return err
		}
	}

	rule := v.ruleCache[value.Type().Name()]
	return v.traverseFields(value, rule)
}

func (v *Validator) traverseFields(value reflect.Value, rule *structRule) error {

	// errors := make([]error, 0)
	for i := 0; i < len(rule.fields); i++ {
		field := value.Field(i)
		for field.Kind() == reflect.Pointer && !field.IsNil() {
			field = field.Elem()
		}

		fieldValue := field.Interface()
		for _, vf := range rule.validateFunc[i] {
			if ok := vf.CheckPass(field.Kind(), fieldValue); !ok {
				// add err
				return ErrorValidateFalse(rule.fields[i].Name, vf.tag)
			}
		}
	}

	return nil
}
