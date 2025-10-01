package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Программные ошибки (ошибки реализации).
var (
	ErrNotStruct        = errors.New("input must be a struct or pointer to struct")
	ErrInvalidValidator = errors.New("invalid validator syntax")
	ErrUnsupportedType  = errors.New("unsupported field type")
	ErrInvalidRegexp    = errors.New("invalid regexp pattern")
	ErrInvalidIntValue  = errors.New("invalid integer value")
	ErrInvalidLenValue  = errors.New("invalid length value")
)

// Ошибки валидации (ошибки входных данных).
var (
	ErrValidationLength = errors.New("validation failed: length")
	ErrValidationRegexp = errors.New("validation failed: regexp")
	ErrValidationIn     = errors.New("validation failed: in")
	ErrValidationMin    = errors.New("validation failed: min")
	ErrValidationMax    = errors.New("validation failed: max")
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

func (v ValidationErrors) Is(target error) bool {
	if target == nil {
		return false
	}

	for _, err := range v {
		if errors.Is(err.Err, target) {
			return true
		}
	}
	return false
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var validationErrors ValidationErrors

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !field.IsExported() {
			continue
		}

		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		if fieldValue.Kind() == reflect.Slice {
			for j := 0; j < fieldValue.Len(); j++ {
				element := fieldValue.Index(j)
				if err := validateField(fmt.Sprintf("%s[%d]", field.Name, j), element, validateTag, &validationErrors); err != nil {
					return err
				}
			}
			continue
		}

		if err := validateField(field.Name, fieldValue, validateTag, &validationErrors); err != nil {
			return err
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateField(
	fieldName string,
	fieldValue reflect.Value,
	validateTag string,
	validationErrors *ValidationErrors,
) error {
	validators := strings.Split(validateTag, "|")
	for _, validator := range validators {
		parts := strings.SplitN(validator, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("%w: %s", ErrInvalidValidator, validator)
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
			return fmt.Errorf("%w: %s", ErrUnsupportedType, fieldValue.Kind())
		}

		if err != nil {
			*validationErrors = append(*validationErrors, ValidationError{
				Field: fieldName,
				Err:   err,
			})
		}
	}

	return nil
}

func validateString(value, validator, arg string) error {
	switch validator {
	case "len":
		expectedLen, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidLenValue, arg)
		}
		if len(value) != expectedLen {
			return ErrValidationLength
		}
	case "regexp":
		re, err := regexp.Compile(arg)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidRegexp, arg)
		}
		if !re.MatchString(value) {
			return ErrValidationRegexp
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
			return ErrValidationIn
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
			return fmt.Errorf("%w: %s", ErrInvalidIntValue, arg)
		}
		if value < minVal {
			return ErrValidationMin
		}
	case "max":
		maxVal, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidIntValue, arg)
		}
		if value > maxVal {
			return ErrValidationMax
		}
	case "in":
		options := strings.Split(arg, ",")
		found := false
		for _, opt := range options {
			optInt, err := strconv.ParseInt(opt, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidIntValue, opt)
			}
			if value == optInt {
				found = true
				break
			}
		}
		if !found {
			return ErrValidationIn
		}
	default:
		return fmt.Errorf("unknown validator %s for int", validator)
	}
	return nil
}
