package babyenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

var (
	// ErrorNotAStructPointer indicates that we were expecting a pointer to a
	// struct but we didn't get it. This is returned when parsing a passed
	// struct.
	ErrorNotAStructPointer = errors.New("expected a pointer to a struct")

	// ErrorUnsupportedType indicates that a struct type isn't supported
	ErrorUnsupportedType = errors.New("unsupported type")
)

// Parse parses a struct for environment variables. We look at the 'env' tag
// for the environment variable names, and the 'envdefault' for the default
// value to the corresponding environment variable.
func Parse(cfg interface{}) error {

	// Make sure we've got a pointer
	val := reflect.ValueOf(cfg)
	if val.Kind() != reflect.Ptr {
		return ErrorNotAStructPointer
	}

	// Make sure our pointer points to a struct
	ref := val.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrorNotAStructPointer
	}

	return parseFields(ref)
}

func parseFields(ref reflect.Value) error {
	for i := 0; i < ref.NumField(); i++ {
		var (
			field     = ref.Field(i)
			fieldKind = ref.Field(i).Kind()
			fieldTags = ref.Type().Field(i).Tag
			fieldName = ref.Type().Field(i).Name
		)

		envVarName := fieldTags.Get("env")
		if len(envVarName) == 0 || envVarName == "-" {
			continue
		}

		if !field.CanSet() {
			return fmt.Errorf("can't set field %v", fieldName)
		}

		envVal := os.Getenv(envVarName)
		defaultVal := fieldTags.Get("envdefault")
		shouldSetDefault := len(envVal) == 0 && len(defaultVal) > 0 && defaultVal != "-"

		switch fieldKind {

		case reflect.String:
			if shouldSetDefault {
				field.SetString(defaultVal)
				continue
			}
			field.SetString(envVal)

		case reflect.Bool:
			if shouldSetDefault {
				if err := setBool(field, defaultVal); err != nil {
					return err
				}
				continue
			}
			if err := setBool(field, envVal); err != nil {
				return err
			}

		case reflect.Int:
			if shouldSetDefault {
				if err := setInt(field, defaultVal); err != nil {
					return err
				}
				continue
			}
			if err := setInt(field, envVal); err != nil {
				return err
			}

		default:
			return ErrorUnsupportedType
		}

	}

	return nil
}

func setBool(v reflect.Value, s string) error {
	if s == "" {
		v.SetBool(false)
		return nil
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

func setInt(v reflect.Value, s string) error {
	if s == "" {
		v.SetInt(0)
		return nil
	}

	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	v.SetInt(n)
	return nil
}
