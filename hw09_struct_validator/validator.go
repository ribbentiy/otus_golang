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

func (e ValidationError) Error() string {
	return e.Err.Error()
}

type ValidationErrors []ValidationError

const tagKey = "validate"

var (
	ErrValidationError     = errors.New("validation error")
	ErrNumberLessThanMin   = errors.New("value less than min")
	ErrNumberMoreThanMax   = errors.New("value more than max")
	ErrValueNotInList      = errors.New("value not in list")
	ErrStringInvalidRegexp = errors.New("value not match regexp pattern")
	ErrStringInvalidLen    = errors.New("value has invalid length")
)

func (v ValidationErrors) Error() string {
	var str strings.Builder
	for _, validationError := range v {
		s := fmt.Sprintf("%v,\n", validationError.Error())
		_, err := str.WriteString(s)
		if err != nil {
			return err.Error()
		}
	}
	return str.String()
}

func Validate(v interface{}) error {
	vType := reflect.TypeOf(v)

	if vType.Kind() != reflect.Struct {
		panic("")
	}

	var validationErrors ValidationErrors

	value := reflect.ValueOf(v)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := vType.Field(i)
		rawTag, ok := vType.Field(i).Tag.Lookup(tagKey)
		if !ok {
			continue
		}
		tags := strings.Split(rawTag, "|")
		for _, tag := range tags {
			rule := strings.Split(tag, ":")
			err := extractedfunc(&fieldType, &field, rule)
			if err != nil {
				if errors.As(err, &ValidationErrors{}) {
					validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
				} else {
					return err
				}
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func extractedfunc(fieldType *reflect.StructField, field *reflect.Value, rule []string) error {
	var validationErrors ValidationErrors
	var verr error
	switch fieldType.Type.Kind() { //nolint:exhaustive
	case reflect.String:
		verr = validateStrings(field.String(), rule)
	case reflect.Int:
		val := int(field.Int())
		verr = validateNumbers(val, rule)
	case reflect.Struct:
		if rule[0] == "nested" {
			err := Validate(field.Interface())
			fmt.Printf("%v: %v \n", fieldType.Name, err)
			if err != nil {
				if errors.As(err, &ValidationErrors{}) {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldType.Name,
						Err:   err,
					})
				} else {
					return err
				}
			}
		}
	case reflect.Slice:
		err := validateSlice(field, fieldType.Name, rule)
		if err != nil {
			if errors.As(err, &ValidationErrors{}) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
			} else {
				return err
			}
		}
	default:
		return errors.New("field type not allowed")
	}

	if verr != nil {
		if errors.Is(verr, ErrValidationError) {
			validationErrors = append(validationErrors, ValidationError{
				Field: fieldType.Name,
				Err:   verr,
			})
		} else {
			return verr
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateStrings(v string, rule []string) error {
	switch rule[0] {
	case "len":
		expectedLen, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		return validateStringLen(v, expectedLen)

	case "regexp":
		return validateStringRegexp(v, rule[1])
	case "in":
		return validateStringInList(v, rule[1])
	default:
		return errors.New("validation method not allowed")
	}
}

func validateNumbers(v int, rule []string) error {
	switch rule[0] {
	case "min":
		min, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		return validateNumberMin(v, min)
	case "max":
		max, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}
		return validateNumberMax(v, max)
	case "in":
		return validateNumberInList(v, rule[1])
	default:
		return errors.New("validation method not allowed")
	}
}

func validateSlice(v *reflect.Value, fieldName string, rule []string) error {
	var validationErrors ValidationErrors
	if !(v.Kind() == reflect.Array || v.Kind() == reflect.Slice) {
		return errors.New("not a slice")
	}
	for i := 0; i < v.Len(); i++ {
		val := v.Index(i)
		var verr error
		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			verr = validateStrings(val.String(), rule)
		case reflect.Int:
			verr = validateNumbers(int(val.Int()), rule)
		default:
			return errors.New("slice values type not allowed")
		}
		if verr != nil {
			if errors.Is(verr, ErrValidationError) {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   verr,
				})
			} else {
				return verr
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateStringLen(v string, expectedLen int) error {
	if len(v) != expectedLen {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidLen)
	}
	return nil
}

func validateStringRegexp(v, pattern string) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	if !re.MatchString(v) {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrStringInvalidRegexp)
	}
	return nil
}

func validateStringInList(v, in string) error {
	if len(in) == 0 {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList)
	}
	list := strings.Split(in, ",")
	for _, s := range list {
		if v == s {
			return nil
		}
	}
	return fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList)
}

func validateNumberMin(v, min int) error {
	if v < min {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrNumberLessThanMin)
	}
	return nil
}

func validateNumberMax(v, max int) error {
	if v > max {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrNumberMoreThanMax)
	}
	return nil
}

func validateNumberInList(v int, in string) error {
	if len(in) == 0 {
		return fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList)
	}
	list := strings.Split(in, ",")
	for _, s := range list {
		listEl, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if v == listEl {
			return nil
		}
	}
	return fmt.Errorf("%w: %v", ErrValidationError, ErrValueNotInList)
}
