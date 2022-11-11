package flago

import (
	"bytes"
	"flag"
	"os"
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
	t.Run("invalid bind: nil-pointer struct", func(t *testing.T) {
		var v *struct{}
		err := Bind(nil, v)
		if err == nil {
			t.Error("error should occur")
		} else if err.Error() != "flago: Bind(nil *struct {})" {
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

func TestBindWithPrefix(t *testing.T) {
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
	err := BindWithPrefix(fs, &v, "pre.")
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	defaults := `  -pre.a int
    	usage of a (default 123)
  -pre.b
    	usage of b (default true)
  -pre.c string
    	usage of c (default "hello world")
  -pre.d value
    	usage of d (default Kim,Machine,Gun)
  -pre.sub.a string
    	usage of sub.a
`
	fs.PrintDefaults()
	if buf.String() != defaults {
		t.Errorf("unexpected defaults: %s", buf.String())
	}

	err = fs.Parse([]string{"-pre.a=456", "-pre.c=Hello World!", "-pre.d=Geon,Kim", "-pre.sub.a=subaval"})
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

func TestBindExpanded(t *testing.T) {
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
	err := BindEnvExpanded(fs, &v)
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	defaults := `  -a value
    	usage of a (default 123)
  -b value
    	usage of b (default true)
  -c value
    	usage of c (default hello world)
  -d value
    	usage of d (default Kim,Machine,Gun)
  -sub.a value
    	usage of sub.a
`
	fs.PrintDefaults()
	if buf.String() != defaults {
		t.Errorf("unexpected defaults: %s", buf.String())
	}

	os.Setenv("FLAGO_A", "678")
	os.Setenv("FLAGO_C", "World Hello!")
	os.Setenv("FLAGO_D", "Kim,Geon")
	os.Setenv("FLAGO_SUB_A", "lavabus")
	err = fs.Parse([]string{"-a=${FLAGO_A}", "-c=${FLAGO_C}", "-d=${FLAGO_D}", "-sub.a=${FLAGO_SUB_A}"})
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	if !reflect.DeepEqual(v, flags{
		A: 678,
		B: true,
		C: "World Hello!",
		D: []string{"Kim", "Geon"},
		e: true,
		Sub: struct {
			A string `flago:"a,usage of sub.a"`
		}{
			A: "lavabus",
		},
	}) {
		t.Errorf("unexpected result: %v", v)
	}
}

func TestBindEnvExpandedWithPrefix(t *testing.T) {
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
	err := BindEnvExpandedWithPrefix(fs, &v, "pre.")
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	defaults := `  -pre.a value
    	usage of a (default 123)
  -pre.b value
    	usage of b (default true)
  -pre.c value
    	usage of c (default hello world)
  -pre.d value
    	usage of d (default Kim,Machine,Gun)
  -pre.sub.a value
    	usage of sub.a
`
	fs.PrintDefaults()
	if buf.String() != defaults {
		t.Errorf("unexpected defaults: %s", buf.String())
	}

	os.Setenv("FLAGO_A", "678")
	os.Setenv("FLAGO_C", "World Hello!")
	os.Setenv("FLAGO_D", "Kim,Geon")
	os.Setenv("FLAGO_SUB_A", "lavabus")
	err = fs.Parse([]string{"-pre.a=${FLAGO_A}", "-pre.c=${FLAGO_C}", "-pre.d=${FLAGO_D}", "-pre.sub.a=${FLAGO_SUB_A}"})
	if err != nil {
		t.Errorf("error should not occur: %v", err)
	}

	if !reflect.DeepEqual(v, flags{
		A: 678,
		B: true,
		C: "World Hello!",
		D: []string{"Kim", "Geon"},
		e: true,
		Sub: struct {
			A string `flago:"a,usage of sub.a"`
		}{
			A: "lavabus",
		},
	}) {
		t.Errorf("unexpected result: %v", v)
	}
}
