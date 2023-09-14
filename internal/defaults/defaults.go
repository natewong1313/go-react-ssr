// https://github.com/meowgorithm/defaults
package defaults

import (
	"fmt"
	"reflect"
	"strconv"
)

// ErrNotAStructPointer indicates that we were expecting a pointer to a struct,
// but got something else.
type ErrNotAStructPointer string

func newErrNotAStructPointer(v interface{}) ErrNotAStructPointer {
	return ErrNotAStructPointer(fmt.Sprintf("%t", v))
}

// Error implements the error interface.
func (e ErrNotAStructPointer) Error() string {
	return fmt.Sprintf("expected a struct, instead got a %T", string(e))
}

// ErrorUnsettable is used when a field cannot be set.
type ErrorUnsettable string

// Error implements the error interface.
func (e ErrorUnsettable) Error() string {
	return fmt.Sprintf("can't set field %s", string(e))
}

// ErrorUnsupportedType indicates that the type of the struct field is not yet
// support in this package.
type ErrorUnsupportedType struct {
	t reflect.Type
}

// Error implements the error interface.
func (e ErrorUnsupportedType) Error() string {
	return fmt.Sprintf("unsupported type %v", e.t)
}

// Apply parses a struct pointer for `default` tags. If the default tag is
// set and the struct member has a default value, the default value will be
// set no the member. Parse expects a struct pointer.
func Apply(t interface{}) error {
	// Make sure we've been given a pointer.
	val := reflect.ValueOf(t)
	if val.Kind() != reflect.Ptr {
		return newErrNotAStructPointer(t)
	}

	// Make sure the pointer is pointing to a struct.
	ref := val.Elem()
	if ref.Kind() != reflect.Struct {
		return newErrNotAStructPointer(t)
	}

	return parseFields(ref)
}

func parseFields(v reflect.Value) error {
	for i := 0; i < v.NumField(); i++ {
		err := parseField(v.Field(i), v.Type().Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func parseField(value reflect.Value, field reflect.StructField) error {
	tagVal := field.Tag.Get("default")

	isStruct := value.Kind() == reflect.Struct
	isStructPointer := value.Kind() == reflect.Ptr && value.Type().Elem().Kind() == reflect.Struct

	if (tagVal == "" || tagVal == "-") && !(isStruct || isStructPointer) {
		return nil
	}

	if !value.CanSet() {
		return ErrorUnsettable(field.Name)
	}

	if !value.IsZero() {
		// A value is set on this field so there's no need to set a default
		// value.
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(tagVal)
		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(tagVal)
		if err != nil {
			return err
		}
		value.SetBool(b)
		return nil

	case reflect.Int:
		i, err := strconv.ParseInt(tagVal, 10, 32)
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil

	case reflect.Int8:
		i, err := strconv.ParseInt(tagVal, 10, 8)
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil

	case reflect.Int16:
		i, err := strconv.ParseInt(tagVal, 10, 16)
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil

	// NB: int32 is also an alias for a rune
	case reflect.Int32:
		i, err := parseInt32(tagVal)
		if err != nil {
			return err
		}
		value.SetInt(int64(i))
		return nil

	case reflect.Int64:
		i, err := strconv.ParseInt(tagVal, 10, 64)
		if err != nil {
			return err
		}
		value.SetInt(i)
		return nil

	case reflect.Uint:
		i, err := strconv.ParseInt(tagVal, 10, 32)
		if err != nil {
			return err
		}
		value.SetUint(uint64(i))
		return nil

	case reflect.Uint8:
		i, err := strconv.ParseInt(tagVal, 10, 8)
		if err != nil {
			return err
		}
		value.SetUint(uint64(i))
		return nil

	case reflect.Uint16:
		i, err := strconv.ParseInt(tagVal, 10, 16)
		if err != nil {
			return err
		}
		value.SetUint(uint64(i))
		return nil

	case reflect.Uint32:
		i, err := strconv.ParseInt(tagVal, 10, 32)
		if err != nil {
			return err
		}
		value.SetUint(uint64(i))
		return nil

	case reflect.Uint64:
		i, err := strconv.ParseInt(tagVal, 10, 64)
		if err != nil {
			return err
		}
		value.SetUint(uint64(i))
		return nil

	case reflect.Float32:
		f, err := strconv.ParseFloat(tagVal, 32)
		if err != nil {
			return err
		}
		value.SetFloat(f)
		return nil

	case reflect.Float64:
		f, err := strconv.ParseFloat(tagVal, 64)
		if err != nil {
			return err
		}
		value.SetFloat(f)
		return nil

	case reflect.Slice:
		switch value.Type().Elem().Kind() {
		// a []uint8 is a an alias for a []byte
		case reflect.Uint8:
			value.SetBytes([]byte(tagVal))
			return nil

		default:
			return ErrorUnsupportedType{value.Type()}
		}

	case reflect.Struct:
		if value.NumField() == 0 {
			return nil
		}
		return parseFields(value) // recurse

	case reflect.Ptr:
		ref := value.Type().Elem()

		switch ref.Kind() {
		case reflect.String:
			value.Set(reflect.ValueOf(&tagVal))
			return nil

		case reflect.Bool:
			b, err := strconv.ParseBool(tagVal)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(&b))
			return nil

		case reflect.Int:
			n, err := strconv.ParseInt(tagVal, 10, 32)
			if err != nil {
				return err
			}
			i := int(n)
			value.Set(reflect.ValueOf(&i))
			return nil

		case reflect.Int8:
			n, err := strconv.ParseInt(tagVal, 10, 8)
			if err != nil {
				return err
			}
			i := int8(n)
			value.Set(reflect.ValueOf(&i))
			return nil

		case reflect.Int16:
			n, err := strconv.ParseInt(tagVal, 10, 16)
			if err != nil {
				return err
			}
			i := int16(n)
			value.Set(reflect.ValueOf(&i))
			return nil

		case reflect.Int32:
			// NB: *int32 is an alias for a *rune
			i, err := parseInt32(tagVal)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(&i))
			return nil

		case reflect.Int64:
			i, err := strconv.ParseInt(tagVal, 10, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(&i))
			return nil

		case reflect.Uint:
			n, err := strconv.ParseInt(tagVal, 10, 32)
			if err != nil {
				return err
			}
			u := uint(n)
			value.Set(reflect.ValueOf(&u))
			return nil

		case reflect.Uint8:
			n, err := strconv.ParseInt(tagVal, 10, 8)
			if err != nil {
				return err
			}
			u := uint8(n)
			value.Set(reflect.ValueOf(&u))
			return nil

		case reflect.Uint16:
			n, err := strconv.ParseInt(tagVal, 10, 16)
			if err != nil {
				return err
			}
			u := uint16(n)
			value.Set(reflect.ValueOf(&u))
			return nil

		case reflect.Uint32:
			n, err := strconv.ParseInt(tagVal, 10, 32)
			if err != nil {
				return err
			}
			u := uint32(n)
			value.Set(reflect.ValueOf(&u))
			return nil

		case reflect.Uint64:
			n, err := strconv.ParseInt(tagVal, 10, 64)
			if err != nil {
				return err
			}
			u := uint64(n)
			value.Set(reflect.ValueOf(&u))
			return nil

		case reflect.Float32:
			f, err := strconv.ParseFloat(tagVal, 32)
			if err != nil {
				return err
			}
			f32 := float32(f)
			value.Set(reflect.ValueOf(&f32))
			return nil

		case reflect.Float64:
			f, err := strconv.ParseFloat(tagVal, 64)
			if err != nil {
				return err
			}
			value.Set(reflect.ValueOf(&f))
			return nil

		case reflect.Slice:
			switch ref.Elem().Kind() {
			// a *[]uint is an alias for *[]byte
			case reflect.Uint8:
				b := []byte(tagVal)
				value.Set(reflect.ValueOf(&b))
				return nil

			default:
				return ErrorUnsupportedType{value.Type()}
			}

		case reflect.Struct:
			if ref.NumField() == 0 {
				return nil
			}

			// If it's nil set it to it's default value so we can set the
			// children if we need to.
			if value.IsNil() {
				value.Set(reflect.New(ref))
			}
			return parseFields(value.Elem()) // recurse

		default:
			return ErrorUnsupportedType{value.Type()}
		}

	default:
		return ErrorUnsupportedType{value.Type()}
	}
}

// Attempt to parse a string as an int32 and, failing that, a rune.
func parseInt32(s string) (int32, error) {
	// Try parsing it as an int.
	i, err := strconv.ParseInt(s, 10, 32)
	if err == nil {
		return int32(i), nil
	}

	// We couldn't parse it as an int, maybe it's a rune.
	runes := []rune(s)
	if len(runes) == 1 {
		return runes[0], nil
	} else {
		return 0, err
	}
}
