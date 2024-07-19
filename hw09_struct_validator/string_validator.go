package hw09structvalidator

import (
	"log"
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
			log.Println("len value is not int")
		}

		stringLength := len(value)

		if stringLength != refValue {
			return ErrLength
		}

	// Check for regex match.
	case "regexp":
		match, err := regexp.MatchString(rule.refValue, value)
		if err != nil {
			log.Println("invalid regexp")
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
		log.Printf("unknown rule name %s", rule.name)
	}

	return nil
}
