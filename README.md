[![PkgGoDev](https://pkg.go.dev/badge/github.com/KimMachineGun/flago)](https://pkg.go.dev/github.com/KimMachineGun/flago)
[![Go Report Card](https://goreportcard.com/badge/github.com/KimMachineGun/flago)](https://goreportcard.com/report/github.com/KimMachineGun/flago)

# Flago

Super simple package for binding command-line flags to your struct.

## Installation

```sh
go get github.com/KimMachineGun/flago
```

## Example

[Playground](https://go.dev/play/p/9RW5vxAFXlh)

```go
package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/KimMachineGun/flago"
)

type Flags struct {
	A string                `flago:"a,usage of a"`
	B int                   `flago:"b,usage of b"`
	C CommaSeparatedStrings `flago:"c,usage of c"`
	D time.Time             `flago:"d,usage of d"`
	E NestedFlags           `flago:"e.,nested flags: "`
}

type NestedFlags struct {
	A string `flago:"a,usage of a"`
	B int    `flago:"b,usage of b"`
}

// CommaSeparatedStrings implements flag.Value.
type CommaSeparatedStrings []string

func (s *CommaSeparatedStrings) String() string {
	return strings.Join(*s, ",")
}

func (s *CommaSeparatedStrings) Set(v string) error {
	*s = strings.Split(v, ",")
	return nil
}

// go run main.go -b=360 -c=Hello,World -d=2020-01-02T00:00:00Z -e.a=CD
func main() {
	// set default values
	flags := Flags{
		A: "AB",
		B: 180,
		C: []string{"Foo", "Bar"},
		D: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		E: NestedFlags{
			A: "AB",
			B: 180,
		},
	}

	// go run main.go --help
	// Output:
	//  -a string
	//        usage of a (default "AB")
	//  -b int
	//        usage of b (default 180)
	//  -c value
	//        usage of c (default Foo,Bar)
	//  -d value
	//        usage of d (default 2020-01-01T00:00:00Z)
	//  -e.a string
	//        nested flags: usage of a (default "AB")
	//  -e.b int
	//        nested flags: usage of b (default 180)

	// parse flags
	err := flago.Parse(&flags)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(flags)
	// Output: {AB 360 [Hello World] 2020-01-02 00:00:00 +0000 UTC {CD 180}}
}
```
