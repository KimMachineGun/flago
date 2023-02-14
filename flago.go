package flago

import (
	"encoding"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

// An InvalidBindError describes an invalid argument passed to Bind.
type InvalidBindError struct {
	Type reflect.Type
}

func (e *InvalidBindError) Error() string {
	if e.Type == nil {
		return "flago: Bind(nil)"
	}
	if e.Type.Kind() != reflect.Ptr || e.Type.Elem().Kind() != reflect.Struct {
		return "flago: Bind(non-struct-pointer " + e.Type.String() + ")"
	}

	return "flago: Bind(nil " + e.Type.String() + ")"
}

// Parse parses the default command-line flag set defined in the flag package.
func Parse(v interface{}) error {
	return ParseWithPrefix(v, "")
}

// ParseWithPrefix parses the default command-line flag set defined in the flag package with prefix.
func ParseWithPrefix(v interface{}, prefix string) error {
	err := BindWithPrefix(flag.CommandLine, v, prefix)
	if err != nil {
		return err
	}
	return flag.CommandLine.Parse(os.Args[1:])
}

// ParseEnvExpanded parses the default command-line flag set defined in the flag package.
// Unlike Parse, it expands environment variables in the flag values.
func ParseEnvExpanded(v interface{}) error {
	return ParseEnvExpandedWithPrefix(v, "")
}

// ParseEnvExpandedWithPrefix parses the default command-line flag set defined in the flag package with prefix.
// Unlike ParseWithPrefix, it expands environment variables in the flag values.
func ParseEnvExpandedWithPrefix(v interface{}, prefix string) error {
	err := BindEnvExpandedWithPrefix(flag.CommandLine, v, prefix)
	if err != nil {
		return err
	}
	return flag.CommandLine.Parse(os.Args[1:])
}

// Bind defines flags based on the struct field tags and
// binds flags to the corresponding fields. If v is nil or
// not a pointer of struct, Bind returns InvalidBindError.
//
// Bind uses 'flago' field tag to specify the name and usage
// of the flag. If a field does not have a 'flago' field tag,
// it will be ignored.
//
//	If a field is struct type, Bind will parse it recursively,
//	and its field tag will be used as a prefix of the names of
//	the flags defined by itself.
//
// Supported Field Types:
//
//   - string
//   - bool
//   - int
//   - int64
//   - uint
//   - uint64
//   - float64
//   - time.Duration
//   - flag.Value
//
// Examples:
//
//	type Example struct {
//		// Name defines a 'name' flag, and its usage message is not set (empty string).
//		Name string `flago:"name"`
//
//		// Age defines a 'age' flag, and its usage message is 'the age of gopher'.
//		// The name and usage specified in the field tag are separated by comma.
//		// flag.IntVar()
//		Age int `flago:"age,the age of gopher"`
//	}
func Bind(fs *flag.FlagSet, v interface{}) error {
	return BindWithPrefix(fs, v, "")
}

// BindWithPrefix defines flags with prefix.
// See the comments of Bind for more details.
func BindWithPrefix(fs *flag.FlagSet, v interface{}, prefix string) error {
	return bind(fs, v, prefix, false)
}

func bind(fs *flag.FlagSet, v interface{}, prefix string, expand bool) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct || rv.IsNil() {
		return &InvalidBindError{Type: reflect.TypeOf(v)}
	}

	elem := reflect.ValueOf(v).Elem()
	elemType := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if !field.CanSet() {
			continue
		}

		rawTag, ok := elemType.Field(i).Tag.Lookup("flago")
		if !ok {
			continue
		}

		tag := strings.SplitN(rawTag, ",", 2)
		if len(tag) != 2 {
			tag = append(tag, "")
		}
		name, usage := prefix+tag[0], tag[1]

		switch f := field.Addr().Interface().(type) {
		case *string:
			if expand {
				fs.Var(newEnvString(*f, f), name, usage)
			} else {
				fs.StringVar(f, name, *f, usage)
			}
		case *bool:
			if expand {
				fs.Var(newEnvBool(*f, f), name, usage)
			} else {
				fs.BoolVar(f, name, *f, usage)
			}
		case *int:
			if expand {
				fs.Var(newEnvInt(*f, f), name, usage)
			} else {
				fs.IntVar(f, name, *f, usage)
			}
		case *int64:
			if expand {
				fs.Var(newEnvInt64(*f, f), name, usage)
			} else {
				fs.Int64Var(f, name, *f, usage)
			}
		case *uint:
			if expand {
				fs.Var(newEnvUint(*f, f), name, usage)
			} else {
				fs.UintVar(f, name, *f, usage)
			}
		case *uint64:
			if expand {
				fs.Var(newEnvUint64(*f, f), name, usage)
			} else {
				fs.Uint64Var(f, name, *f, usage)
			}
		case *float64:
			if expand {
				fs.Var(newEnvFloat64(*f, f), name, usage)
			} else {
				fs.Float64Var(f, name, *f, usage)
			}
		case *time.Duration:
			if expand {
				fs.Var(newEnvDuration(*f, f), name, usage)
			} else {
				fs.DurationVar(f, name, *f, usage)
			}
		case flag.Value:
			if expand {
				f = newEnvVar(f)
			}
			fs.Var(f, name, usage)
		case interface {
			encoding.TextMarshaler
			encoding.TextUnmarshaler
		}:
			if expand {
				f = newEnvText(f)
			}
			fs.TextVar(f, name, f, usage)
		default:
			if field.Kind() != reflect.Struct {
				return fmt.Errorf("unsupported type: %T", f)
			}

			if err := bind(fs, f, name, expand); err != nil {
				return err
			}
		}
	}

	return nil
}

// BindEnvExpanded defines env-expanded flags based on the struct field tags and
// binds flags to the corresponding fields.
func BindEnvExpanded(fs *flag.FlagSet, v interface{}) error {
	return BindEnvExpandedWithPrefix(fs, v, "")
}

// BindEnvExpandedWithPrefix defines env-expanded flags with prefix.
func BindEnvExpandedWithPrefix(fs *flag.FlagSet, v interface{}, prefix string) error {
	return bind(fs, v, prefix, true)
}
