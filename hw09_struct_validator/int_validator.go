package hw09structvalidator

import (
	"log"
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
			log.Println("reference value is not int")
		}

		if value < refValue {
			return ErrMin
		}

	// Check for maximum.
	case "max":
		// Get rule reference value.
		refValue, err := strconv.Atoi(rule.refValue)
		if err != nil {
			log.Println("reference value is not int")
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
				log.Println("reference value is not int")
				continue
			}

			if value == refValue {
				return nil
			}

			return ErrNotIn
		}

	default:
		log.Printf("unknown rule name %s", rule.name)
	}

	return nil
}
