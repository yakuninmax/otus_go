package hw09structvalidator

import (
	"errors"
	"reflect"
	"strings"
)

const (
	tagName        = "validate"
	rulesDelimiter = "|"
	rangeDelimeter = ","
	delimeter      = ":"
)

var (
	ErrNotStruct     = errors.New("value is not a structure")
	ErrMin           = errors.New("value is less than minimum")
	ErrMax           = errors.New("value is greater than maximum")
	ErrLength        = errors.New("invalid string length")
	ErrRegexp        = errors.New("not match regexp")
	ErrInvalidRegexp = errors.New("invalid regexp")
	ErrNotIn         = errors.New("not in range")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

// Type for rules.
type rule struct {
	name     string
	refValue string
}

func (v ValidationErrors) Error() string {
	str := strings.Builder{}

	for _, errStruct := range v {
		str.WriteString(errStruct.Err.Error())
	}

	return str.String()
}

func Validate(v interface{}) error {
	// Check input struct.
	inputStruct := reflect.ValueOf(v)
	if inputStruct.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	// Create validation errors slice.
	var validationErrors ValidationErrors

	inputStructType := inputStruct.Type()

	for i := 0; i < inputStructType.NumField(); i++ {
		// Get field.
		fieldValue := inputStruct.Field(i)

		// Get rules.
		rules := getRules(inputStructType.Field(i).Tag.Get(tagName))

		// Validate field.
		err := validateField(fieldValue, rules)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{Field: inputStructType.Field(i).Name, Err: err})
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

// Get validation rules from tag.
func getRules(fieldTag string) []rule {
	rawRules := strings.Split(fieldTag, rulesDelimiter)

	// Create rules slice.
	rules := make([]rule, 0)

	// Get rules.
	for _, rawRule := range rawRules {
		// Get rule values.
		values := strings.Split(rawRule, delimeter)

		// Check rule values count.
		if len(values) != 2 {
			continue
		}

		// Create rule structure.
		rule := rule{
			name:     values[0],
			refValue: values[1],
		}

		// Append rule to slice.
		rules = append(rules, rule)
	}

	return rules
}

// Validate field.
func validateField(fieldValue reflect.Value, rules []rule) error {
	var err error

	// Get field value as interface.
	fieldValueType := fieldValue.Interface()

	// Switch by interface type.
	switch fieldValueType.(type) {
	// If value is integer.
	case int:
		for _, rule := range rules {
			err = intValidate(fieldValue, rule)
			if err != nil {
				return err
			}
		}

	// If value is string.
	case string:
		for _, rule := range rules {
			err = stringValidate(fieldValue, rule)
			if err != nil {
				return err
			}
		}

	// If value is slice.
	case []int, []string:
		for i := 0; i < fieldValue.Len(); i++ {
			value := fieldValue.Index(i)
			err = validateField(value, rules)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
