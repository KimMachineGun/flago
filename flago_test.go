package flago

import (
	"bytes"
	"flag"
	"reflect"
	"strings"
	"testing"
)

func TestBind__InvalidBindError(t *testing.T) {
	t.Run("invalid bind: nil", func(t *testing.T) {
		err := Bind(nil, nil)
		if err == nil {
			t.Error("error should occur")
		} else if err.Error() != "flago: Bind(nil)" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("invalid bind: non-struct", func(t *testing.T) {
		var v int
		err := Bind(nil, &v)
		if err == nil {
			t.Error("error should occur")
		} else if err.Error() != "flago: Bind(non-struct-pointer *int)" {
			t.Errorf("unexpected error: %v", err)
		}
	})
	t.Run("invalid bind: non-pointer", func(t *testing.T) {
		var v struct{}
		err := Bind(nil, v)
		if err == nil {
			t.Error("error should occur")
		} else if err.Error() != "flago: Bind(non-struct-pointer struct {})" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

type commaSeparated []string

func (f *commaSeparated) String() string {
	return strings.Join(*f, ",")
}

func (f *commaSeparated) Set(s string) error {
	*f = strings.Split(s, ",")
	return nil
}

func TestBind(t *testing.T) {
	fs := flag.NewFlagSet("flago", flag.ContinueOnError)

	buf := bytes.NewBuffer(nil)
	fs.SetOutput(buf)

	type flags struct {
		A   int            `flago:"a,usage of a"`
		B   bool           `flago:"b,usage of b"`
		C   string         `flago:"c,usage of c"`
		D   commaSeparated `flago:"d,usage of d"`
		Sub struct {
			A string `flago:"a,usage of sub.a"`
		} `flago:"sub."`
		// e will be omitted, since it is an unexported field.
		e bool `flago:"e,usage of e"`
	}
	v := flags{
		A: 123,
		B: true,
		C: "hello world",
		D: []string{"Kim", "Machine", "Gun"},
		e: true,
	}
	err := Bind(fs, &v)
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	defaults := `  -a int
    	usage of a (default 123)
  -b	usage of b (default true)
  -c string
    	usage of c (default "hello world")
  -d value
    	usage of d (default Kim,Machine,Gun)
  -sub.a string
    	usage of sub.a
`
	fs.PrintDefaults()
	if buf.String() != defaults {
		t.Errorf("unexpected defaults: %s", buf.String())
	}

	err = fs.Parse([]string{"-a=456", "-c=Hello World!", "-d=Geon,Kim", "-sub.a=subaval"})
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	if !reflect.DeepEqual(v, flags{
		A: 456,
		B: true,
		C: "Hello World!",
		D: []string{"Geon", "Kim"},
		e: true,
		Sub: struct {
			A string `flago:"a,usage of sub.a"`
		}{
			A: "subaval",
		},
	}) {
		t.Errorf("unexpected result: %v", v)
	}
}
