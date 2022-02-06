package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ErrUnsupportedType is returned when an environment variable target to bind is
// not supported. For example a custom struct type.
var ErrUnsupportedType = errors.New("unsupported type")

// ErrParse is used in the ParseError type when checking error.Is to be able to
// identify the error responses from the bind functions.
var ErrParse = errors.New("failed to parse")

// ParseError is an error type that unwraps to the internal parsing error
// obtained.
type ParseError struct {
	EnvKey   string
	EnvValue string
	Err      error
}

// Error returns the error string. Makes it compliant with the error interface.
func (err ParseError) Error() string {
	return fmt.Sprintf("env %q=%q: %s", err.EnvKey, err.EnvValue, err.Err)
}

// Is returns true if the target error is ErrParse or if the inner actual
// parsing error is the same error.
//
// This method provides compatibility with the errors.Is function.
func (err ParseError) Is(target error) bool {
	return errors.Is(ErrParse, target) || errors.Is(err.Err, target)
}

// As returns true if this error could be unwrapped into the type of the target
// error.
//
// This method provides compatibility with the errors.As function.
func (err ParseError) As(target any) bool {
	return errors.As(err.Err, target)
}

// Unwrap returns the inner actual parsing error, such as the error return value
// from a strconv function.
//
// This method provides compatibility with the errors.Unwrap function.
func (err ParseError) Unwrap() error {
	return err.Err
}

// BindConstraint is a generic type constraint of all the types that the Bind
// function supports.
type BindConstraint interface {
	*string | *bool | *int | *int32 | *int64 | *uint | *uint32 | *uint64 |
		*float32 | *float64 | *time.Time | *time.Duration
}

// Bind will take a value pointer and depending on its type will try to parse
// the environment variable, if set and not empty, using the appropriate parsing
// function.
//
// If the environment variable is not set, is empty, or the function returns an
// error, the value of the target interface is left unchanged.
//
// Returns an env.ParseError on parsing errors.
//
// Returns a wrapped env.ErrUnsupportedType error if the type of the interface
// is not supported.
//
// Returns a wrapped env.ErrNotAPointer error if the target interface is not a
// pointer.
//
// Returns nil otherwise.
func Bind[T BindConstraint](i T, key string) error {
	var envStr, ok = LookupNoEmpty(key)
	if !ok {
		return nil
	}
	switch ptr := (any)(i).(type) {
	case *string:
		*ptr = envStr
	case *bool:
		value, err := strconv.ParseBool(envStr)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	case *int:
		value, err := strconv.ParseInt(envStr, 10, strconv.IntSize)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = int(value)
	case *int32:
		value, err := strconv.ParseInt(envStr, 10, 32)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = int32(value)
	case *int64:
		value, err := strconv.ParseInt(envStr, 10, 64)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	case *uint:
		value, err := strconv.ParseUint(envStr, 10, strconv.IntSize)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = uint(value)
	case *uint32:
		value, err := strconv.ParseUint(envStr, 10, 32)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = uint32(value)
	case *uint64:
		value, err := strconv.ParseUint(envStr, 10, 64)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	case *float32:
		value, err := strconv.ParseFloat(envStr, 32)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = float32(value)
	case *float64:
		value, err := strconv.ParseFloat(envStr, 64)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	case *time.Time:
		value, err := time.Parse(time.RFC3339, envStr)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	case *time.Duration:
		value, err := time.ParseDuration(envStr)
		if err != nil {
			return ParseError{key, envStr, err}
		}
		*ptr = value
	default:
		return fmt.Errorf("env %q: %w: %T", key, ErrUnsupportedType, i)
	}
	return nil
}

// BindMultiple updates the Go variables via the pointers with the values of the
// environment variables, if set and not empty, for each respective pair in
// the map.
//
// If the environment variable is not set, is empty, or the function returns an
// error, the value of the respective target interface is left unchanged.
//
// An error is returned if any of the bindings failed to bind.
func BindMultiple[T BindConstraint](bindings map[T]string) error {
	for ptr, key := range bindings {
		if err := Bind(ptr, key); err != nil {
			return err
		}
	}
	return nil
}

// LookupNoEmpty retrieves the value of the environment variable.
//
// Returns ("", false) if the environment variable was empty, or not set.
// Returns (envVariableValue, true) otherwise.
func LookupNoEmpty(key string) (string, bool) {
	var str, ok = os.LookupEnv(key)
	if str == "" {
		return "", false
	}
	return str, ok
}
