package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range v {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(fmt.Sprintf("%s: %v", err.Field, err.Err))
	}
	return sb.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("input must be a struct or pointer to struct")
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if field.PkgPath != "" {
			continue
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Handle slices
		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < fieldValue.Len(); j++ {
				element := fieldValue.Index(j)
				if errs := validateField(fmt.Sprintf("%s[%d]", field.Name, j), element, validateTag); errs != nil {
					validationErrors = append(validationErrors, errs...)
				}
			}
			continue
		}

		if errs := validateField(field.Name, fieldValue, validateTag); errs != nil {
			validationErrors = append(validationErrors, errs...)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(fieldName string, fieldValue reflect.Value, validateTag string) ValidationErrors {
	var validationErrors ValidationErrors

	validators := strings.Split(validateTag, "|")
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("invalid validator syntax: %s", validator),
			})
			continue
		}

		validatorName := parts[0]
		validatorValue := parts[1]

		var err error
		//nolint:exhaustive
		switch fieldValue.Kind() {
		case reflect.String:
			err = validateString(fieldValue.String(), validatorName, validatorValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = validateInt(fieldValue.Int(), validatorName, validatorValue)
		default:
			err = fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
		}

		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	return validationErrors
}

func validateString(value, validator, arg string) error {
	switch validator {
	case "len":
		expectedLen, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid length value: %s", arg)
		}
		if len(value) != expectedLen {
			return fmt.Errorf("length must be %d", expectedLen)
		}
	case "regexp":
		matched, err := regexp.MatchString(arg, value)
		if err != nil {
			return fmt.Errorf("invalid regexp: %s", arg)
		}
		if !matched {
			return fmt.Errorf("must match regexp %s", arg)
		}
	case "in":
		options := strings.Split(arg, ",")
		found := false
		for _, opt := range options {
			if value == opt {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("must be one of %v", options)
		}
	default:
		return fmt.Errorf("unknown validator %s for string", validator)
	}
	return nil
}

func validateInt(value int64, validator, arg string) error {
	switch validator {
	case "min":
		minVal, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid min value: %s", arg)
		}
		if value < minVal {
			return fmt.Errorf("must be >= %d", minVal)
		}
	case "max":
		maxVal, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid max value: %s", arg)
		}
		if value > maxVal {
			return fmt.Errorf("must be <= %d", maxVal)
		}
	case "in":
		options := strings.Split(arg, ",")
		found := false
		for _, opt := range options {
			optInt, err := strconv.ParseInt(opt, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid in value: %s", opt)
			}
			if value == optInt {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("must be one of %v", options)
		}
	default:
		return fmt.Errorf("unknown validator %s for int", validator)
	}
	return nil
}
