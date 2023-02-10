package flago

import (
	"encoding"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type envString string

func newEnvString(val string, p *string) *envString {
	*p = val
	return (*envString)(p)
}

func (p *envString) Set(s string) error {
	*p = envString(os.ExpandEnv(s))
	return nil
}

func (p *envString) String() string {
	return string(*p)
}

type envBool bool

func newEnvBool(val bool, p *bool) *envBool {
	*p = val
	return (*envBool)(p)
}

func (p *envBool) Set(s string) error {
	v, err := strconv.ParseBool(os.ExpandEnv(s))
	if err != nil {
		return fmt.Errorf("invalid bool value: %v", v)
	}
	*p = envBool(v)
	return nil
}

func (p *envBool) String() string {
	return strconv.FormatBool(bool(*p))
}

type envInt int

func newEnvInt(val int, p *int) *envInt {
	*p = val
	return (*envInt)(p)
}

func newEnvIntValue(val int, p *int) *envInt {
	*p = val
	return (*envInt)(p)
}

func (p *envInt) Set(s string) error {
	v, err := strconv.ParseInt(os.ExpandEnv(s), 0, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("invalid int value: %v", v)
	}
	*p = envInt(v)
	return nil
}

func (p *envInt) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

type envInt64 int64

func newEnvInt64(val int64, p *int64) *envInt64 {
	*p = val
	return (*envInt64)(p)
}

func (p *envInt64) Set(s string) error {
	v, err := strconv.ParseInt(os.ExpandEnv(s), 0, 64)
	if err != nil {
		return fmt.Errorf("invalid int64 value: %v", v)
	}
	*p = envInt64(v)
	return nil
}

func (p *envInt64) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

type envUint uint

func newEnvUint(val uint, p *uint) *envUint {
	*p = val
	return (*envUint)(p)
}

func (p *envUint) Set(s string) error {
	v, err := strconv.ParseUint(os.ExpandEnv(s), 0, strconv.IntSize)
	if err != nil {
		return fmt.Errorf("invalid uint value: %v", v)
	}
	*p = envUint(v)
	return nil
}

func (p *envUint) String() string {
	return strconv.FormatUint(uint64(*p), 10)
}

type envUint64 uint64

func newEnvUint64(val uint64, p *uint64) *envUint64 {
	*p = val
	return (*envUint64)(p)
}

func (p *envUint64) Set(s string) error {
	v, err := strconv.ParseUint(os.ExpandEnv(s), 0, 64)
	if err != nil {
		return fmt.Errorf("invalid uint64 value: %v", v)
	}
	*p = envUint64(v)
	return nil
}

func (p *envUint64) String() string {
	return strconv.FormatUint(uint64(*p), 10)
}

type envFloat64 float64

func newEnvFloat64(val float64, p *float64) *envFloat64 {
	*p = val
	return (*envFloat64)(p)
}

func (p *envFloat64) Set(s string) error {
	v, err := strconv.ParseFloat(os.ExpandEnv(s), 64)
	if err != nil {
		return fmt.Errorf("invalid float64 value: %v", v)
	}
	*p = envFloat64(v)
	return nil
}

func (p *envFloat64) String() string {
	return strconv.FormatFloat(float64(*p), 'g', -1, 64)
}

type envDuration time.Duration

func newEnvDuration(val time.Duration, p *time.Duration) *envDuration {
	*p = val
	return (*envDuration)(p)
}

func (p *envDuration) Set(s string) error {
	v, err := time.ParseDuration(os.ExpandEnv(s))
	if err != nil {
		return fmt.Errorf("invalid duration value: %v", v)
	}
	*p = envDuration(v)
	return nil
}

func (p *envDuration) String() string {
	return time.Duration(*p).String()
}

type envVar struct {
	v flag.Value
}

func newEnvVar(val flag.Value) *envVar {
	return &envVar{val}
}

func (p *envVar) Set(s string) error {
	return p.v.Set(os.ExpandEnv(s))
}

func (p *envVar) String() string {
	if p.v == nil {
		return ""
	}
	return p.v.String()
}

type envText struct {
	v interface {
		encoding.TextMarshaler
		encoding.TextUnmarshaler
	}
}

func newEnvText(val interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}) *envText {
	return &envText{val}
}

func (p *envText) MarshalText() ([]byte, error) {
	return p.v.MarshalText()
}

func (p *envText) UnmarshalText(b []byte) error {
	if len(b) > 0 {
		b = []byte(os.ExpandEnv(string(b)))
	}
	return p.v.UnmarshalText(b)
}
