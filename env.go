// Package babyenv collects environment variables and places them in
// corresponding struct fields. It aims to reduce the boilerplate in reading
// data from the environment.
//
// The struct should contain `env` tags indicating the names of corresponding
// environment variables. The values of those environment variables will be
// then collected and placed into the struct. If nothing is found, struct
// fields will be given their default values (for example, `bool`s will be
// `false`).
//
//     type config struct {
//         Name string `env:"NAME"`
//     }
//
// Default values can also be provided in the `default` tag.
//
//     `env:"NAME" default:"Jane"`
//
// A 'required' flag can also be set in the following format:
//
//     `env:"NAME,required"`
//
// If a required flag is set the 'default' tag will be ignored.
//
// Only a few types are supported: string, bool, int, []byte, *string, *bool,
// *int, *[]byte. An error will be returned if other types are attempted to
// be processed.
//
// Example:
//
//     package main
//
//     import (
//         "fmt"
//         "os"
//         "github.com/magicnumbers/babyenv"
//     )
//
//     type config struct {
//         Debug   bool   `env:"DEBUG"`
//         Port    string `env:"PORT" default:"8000"`
//         Workers int    `env:"WORKERS" default:"16"`
//         Name    string `env:"NAME,required"`
//     }
//
//     func main() {
//         os.Setenv("DEBUG", "true")
//         os.Setenv("WORKERS", "4")
//         os.Setenv("NAME", "Jane")
//
//         var cfg config
//         if err := babyenv.Parse(&cfg); err != nil {
//             log.Fatalf("could not get environment vars: %v", err)
//         }
//
//         fmt.Printf("%b\n%s\n%d\n%s", cfg.Debug, cfg.Port, cfg.Workers, cfg.Name)
//
//         // Output:
//         // true
//         // 8000
//         // 4
//         // Jane
//     }
package babyenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	// ErrorNotAStructPointer indicates that we were expecting a pointer to a
	// struct but we didn't get it. This is returned when parsing a passed
	// struct.
	ErrorNotAStructPointer = errors.New("expected a pointer to a struct")
)

// ErrorUnsettable is used when a field cannot be set
type ErrorUnsettable struct {
	FieldName string
}

// Error implements the error interface
func (e *ErrorUnsettable) Error() string {
	return fmt.Sprintf("can't set field %s", e.FieldName)
}

// ErrorUnsupportedType is used when we attempt to parse a struct field of an
// unsupported type
type ErrorUnsupportedType struct {
	Type reflect.Type
}

// Error implements the error interface
func (e *ErrorUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type %v", e.Type)
}

// ErrorEnvVarRequired is used when a `required` flag is used and the value of
// the corresponding environment variable is empty
type ErrorEnvVarRequired struct {
	Name string
}

// Error implements the error interface
func (e *ErrorEnvVarRequired) Error() string {
	return fmt.Sprintf("%s is required", e.Name)
}

// Parse parses a struct for environment variables, placing found values in the
// struct, altering it. We look at the 'env' tag for the environment variable
// names, and the 'default' for the default value to the corresponding
// environment variable.
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

// Interate over the fields of a struct, looking for `env` tags indicating
// environment variable names and `default` inicating default values. We're
// expecting a pointer to a struct here, and either environment variables or
// defaults will be placed in the struct. If a non-struct pointer is passed we
// return an error.
//
// Note that a required flag can also be passed in the form of:
//
//     VarName string `env:"VAR_NAME,required"`
//
// If a required flag is set, and the environment variable is empty, the
// `default` tag is ignored.
func parseFields(ref reflect.Value) error {
	for i := 0; i < ref.NumField(); i++ {
		var (
			field      = ref.Field(i)
			fieldKind  = ref.Field(i).Kind()
			fieldTags  = ref.Type().Field(i).Tag
			fieldName  = ref.Type().Field(i).Name
			envVarName string
			required   bool
		)

		tagVal := fieldTags.Get("env")
		if tagVal == "" || tagVal == "-" {
			continue
		}

		if !field.CanSet() {
			return &ErrorUnsettable{fieldName}
		}

		// The tag we're looking at will look something like one of these:
		//
		//     `env:"NAME"`
		//     `env:"NAME,required"`
		//
		// Here we split on the comma and sort out the parts.
		tagValParts := strings.Split(tagVal, ",")
		if len(tagValParts) == 0 { // This should never happen
			continue
		} else if len(tagValParts) >= 1 {
			envVarName = tagValParts[0]
		}
		if len(tagValParts) >= 2 && strings.TrimSpace(tagValParts[1]) == "required" {
			required = true
		}

		// Get the value of the environment var
		envVarVal := os.Getenv(envVarName)

		// Return an error if the required flag is set and the env var is empty
		if envVarVal == "" && required {
			return &ErrorEnvVarRequired{envVarName}
		}

		defaultVal := fieldTags.Get("default")

		// Is the situation such that we should set a default value? We only
		// do it if the value of the given environment varaiable is empty, and
		// we have a non-empty default value.
		shouldSetDefault := len(envVarVal) == 0 && len(defaultVal) > 0 && defaultVal != "-"

		// Set the field accoring to it's kind
		switch fieldKind {

		case reflect.String:
			if shouldSetDefault {
				field.SetString(defaultVal)
				continue
			}
			field.SetString(envVarVal)

		case reflect.Bool:
			if shouldSetDefault {
				if err := setBool(field, defaultVal); err != nil {
					return err
				}
				continue
			}
			if err := setBool(field, envVarVal); err != nil {
				return err
			}

		case reflect.Int:
			if shouldSetDefault {
				if err := setInt(field, defaultVal); err != nil {
					return err
				}
				continue
			}
			if err := setInt(field, envVarVal); err != nil {
				return err
			}

		case reflect.Int64:
			if shouldSetDefault {
				if err := setInt64(field, defaultVal); err != nil {
					return err
				}
				continue
			}
			if err := setInt64(field, envVarVal); err != nil {
				return err
			}

		// Slices are a whole can of worms
		case reflect.Slice:
			switch field.Type().Elem().Kind() {

			// []uint8 is an alias for []byte
			case reflect.Uint8:
				if shouldSetDefault {
					field.SetBytes([]byte(defaultVal))
					continue
				}
				field.SetBytes([]byte(envVarVal))

			default:
				return &ErrorUnsupportedType{field.Type()}

			}

		// Pointers are also a whole other can of worms
		case reflect.Ptr:
			ptr := field.Type().Elem()

			switch ptr.Kind() {

			case reflect.String:
				if shouldSetDefault {
					field.Set(reflect.ValueOf(&defaultVal))
					continue
				}
				field.Set(reflect.ValueOf(&envVarVal))

			case reflect.Bool:
				if shouldSetDefault {
					if err := setBoolPointer(field, defaultVal); err != nil {
						return err
					}
					continue
				}
				if err := setBoolPointer(field, envVarVal); err != nil {
					return err
				}

			case reflect.Int:
				if shouldSetDefault {
					if err := setIntPointer(field, defaultVal); err != nil {
						return err
					}
					continue
				}
				if err := setIntPointer(field, envVarVal); err != nil {
					return err
				}

			case reflect.Int64:
				if shouldSetDefault {
					if err := setInt64Pointer(field, defaultVal); err != nil {
						return err
					}
					continue
				}
				if err := setInt64Pointer(field, envVarVal); err != nil {
					return err
				}

			// A poiner to a slice!! Whole other level
			case reflect.Slice:

				switch ptr.Elem().Kind() {

				// *[]uint8 is an alias for *[]byte
				case reflect.Uint8:
					var byteSlice []byte
					if shouldSetDefault {
						byteSlice = []byte(defaultVal)
					} else {
						byteSlice = []byte(envVarVal)
					}
					field.Set(reflect.ValueOf(&byteSlice))

				default:
					return &ErrorUnsupportedType{field.Type()}

				}

			default:
				return &ErrorUnsupportedType{field.Type()}
			}

		default:
			return &ErrorUnsupportedType{field.Type()}
		}

	}

	return nil
}

func setBool(v reflect.Value, s string) error {
	if s == "" {
		// Default to false
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
		// Default to 0
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

func setInt64(v reflect.Value, s string) error {
	if s == "" {
		// Default to 0
		v.SetInt(0)
		return nil
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(n)
	return nil
}

func setBoolPointer(v reflect.Value, s string) error {
	if s == "" {
		// Default to false
		b := false
		v.Set(reflect.ValueOf(&b))
		return nil
	}

	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(&b))
	return nil
}

func setIntPointer(v reflect.Value, s string) error {
	if s == "" {
		// Default to 0
		n := 0
		v.Set(reflect.ValueOf(&n))
		return nil
	}

	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	i := int(i64)

	v.Set(reflect.ValueOf(&i))
	return nil
}

func setInt64Pointer(v reflect.Value, s string) error {
	if s == "" {
		// Default to 0
		n := 0
		v.Set(reflect.ValueOf(&n))
		return nil
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(&i))
	return nil
}
