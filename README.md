[![PkgGoDev](https://pkg.go.dev/badge/github.com/KimMachineGun/flago)](https://pkg.go.dev/github.com/KimMachineGun/flago)
[![Go Report Card](https://goreportcard.com/badge/github.com/KimMachineGun/flago)](https://goreportcard.com/report/github.com/KimMachineGun/flago)
# Flago
Super simple package for binding command-line flags to your struct.

## Installation
```
go get github.com/KimMachineGun/flago
```

## Example
[Playground](https://play.golang.org/p/c3HlUPZj1Ot)
```go
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/KimMachineGun/flago"
)

type Flags struct {
	A string                `flago:"a,usage of a"`
	B int                   `flago:"b,usage of b"`
	C CommaSeparatedStrings `flago:"c,usage of c"`
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

// go run main.go -b=360 -c=Hello,World
func main() {
	// set default values
	flags := Flags{
		A: "AB",
		B: 180,
		C: []string{"Foo", "Bar"},
	}

	err := flago.Bind(flag.CommandLine, &flags)
	if err != nil {
		log.Fatalln(err)
	}

	flag.PrintDefaults()
	// Output:
	//  -a string
	//        usage of a (default "AB")
	//  -b int
	//        usage of b (default 180)
	//  -c value
	//        usage of c (default Foo,Bar)

	flag.Parse()

	fmt.Println(flags)
	// Output: {AB 360 [Hello World]}
}
```
