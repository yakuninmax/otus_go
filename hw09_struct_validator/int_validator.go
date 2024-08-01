package hw09structvalidator

import (
	"reflect"
	"strconv"
	"strings"
)

// Validate integer values.
func intValidate(fieldValue reflect.Value, rule rule) error {
	// Get fieldValue as int.
	value := int(fieldValue.Int())

	switch rule.name {
	// Check for minimum.
	case "min":
		// Get rule reference value.
		refValue, err := strconv.Atoi(rule.refValue)
		if err != nil {
			return CommonError{CommonErr: err}
		}

		if value < refValue {
			return ErrMin
		}

	// Check for maximum.
	case "max":
		// Get rule reference value.
		refValue, err := strconv.Atoi(rule.refValue)
		if err != nil {
			return CommonError{CommonErr: err}
		}

		if value > refValue {
			return ErrMax
		}

	// Check value in range.
	case "in":
		// Get list of acceptable values.
		acceptableValues := strings.Split(rule.refValue, rangeDelimeter)

		for _, acceptableValue := range acceptableValues {
			// Get reference value.
			refValue, err := strconv.Atoi(acceptableValue)
			if err != nil {
				return CommonError{CommonErr: err}
			}

			if value == refValue {
				return nil
			}
		}

		return ErrNotIn

	default:
		return CommonError{CommonErr: ErrInvalidRule}
	}

	return nil
}
