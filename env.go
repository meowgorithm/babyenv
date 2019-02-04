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
// Default values can also be provided in the `default` tag.
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
//         Debug bool  `env:"DEBUG"`
//         Port string `env:"PORT" default:"8000"`
//         Workers int `env:"WORKERS" default:"16"`
//     }
//
//     func main() {
//         os.Setenv("DEBUG", "true")
//         os.Setenv("WORKERS", "4")
//
//         var cfg config
//         if err := babyenv.Parse(&cfg); err != nil {
//             log.Fatalf("could not get environment vars: %v", err)
//         }
//
//         fmt.Printf("%b\n%s\n%d", cfg.Debug, cfg.Port, cfg.Workers)
//
//         // Output:
//         // true
//         // 8000
//         // 4
//     }
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
)

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

		envVarVal := os.Getenv(envVarName)
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

		// Pointers are a whole other can of worms
		case reflect.Ptr:
			switch field.Type().Elem().Kind() {

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

			default:
				return fmt.Errorf("unsupported type %v", field.Type())
			}

		default:
			return fmt.Errorf("unsupported type %v", field.Type())
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
