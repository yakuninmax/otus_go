package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Validate string values.
func stringValidate(fieldValue reflect.Value, rule rule) error {
	// Get field as sring value.
	value := fieldValue.String()

	switch rule.name {
	// Check string length.
	case "len":
		refValue, err := strconv.Atoi(rule.refValue)
		if err != nil {
			return CommonError{CommonErr: ErrInvalidRef}
		}

		stringLength := len(value)

		if stringLength != refValue {
			return ErrLength
		}

	// Check for regex match.
	case "regexp":
		match, err := regexp.MatchString(rule.refValue, value)
		if err != nil {
			return CommonError{CommonErr: ErrInvalidRegexp}
		}

		if !match {
			return ErrRegexp
		}

	// Check value in range.
	case "in":
		acceptableValues := strings.Split(rule.refValue, rangeDelimeter)

		for _, acceptableValue := range acceptableValues {
			if value == acceptableValue {
				return nil
			}
		}

		return ErrNotIn

	default:
		return CommonError{CommonErr: ErrInvalidRule}
	}

	return nil
}
