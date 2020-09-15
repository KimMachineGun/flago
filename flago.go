package flago

import (
	"flag"
	"fmt"
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

// Bind defines flags based on the struct field tags and
// binds flags to the corresponding fields.
//
// The name of field tag is 'flago' and its value is used to
// specify the name and usage of the flag.
//
// Supported Field Types:
//
//  - string
//  - bool
//  - int
//  - int64
//  - uint
//  - uint64
//  - float64
//  - time.Duration
//  - flag.Value
//
// Examples:
//
//   // Name defines a 'name' flag, and its usage is skipped.
//   Name string `flago:"name"`
//
//   // Age defines a 'name' flag, and its usage is 'the age of gopher'.
//   // The name and usage specified in field tag are separated by comma.
//   Age int `flago:"age,the age of gopher"`
//
func Bind(fs *flag.FlagSet, v interface{}) error {
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

		rawTag := elemType.Field(i).Tag.Get("flago")
		if rawTag == "" {
			continue
		}

		tag := strings.SplitN(rawTag, ",", 2)
		if len(tag) != 2 {
			tag = append(tag, "")
		}
		name, usage := strings.TrimSpace(tag[0]), strings.TrimSpace(tag[1])

		switch f := field.Addr().Interface().(type) {
		case *string:
			fs.StringVar(f, name, *f, usage)
		case *bool:
			fs.BoolVar(f, name, *f, usage)
		case *int:
			fs.IntVar(f, name, *f, usage)
		case *int64:
			fs.Int64Var(f, name, *f, usage)
		case *uint:
			fs.UintVar(f, name, *f, usage)
		case *uint64:
			fs.Uint64Var(f, name, *f, usage)
		case *float64:
			fs.Float64Var(f, name, *f, usage)
		case *time.Duration:
			fs.DurationVar(f, name, *f, usage)
		case flag.Value:
			fs.Var(f, name, usage)
		default:
			return fmt.Errorf("unsupported type: %T", f)
		}
	}

	return nil
}
