package fixedwidth

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	errInvalidDateFormat = fmt.Errorf("invalid date format")
)

type encoder struct {
	bytes.Buffer
}

func (e encoder) encode(v interface{}) ([]byte, error) {
	value := reflect.ValueOf(v)

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	valueType := value.Type()

	for i := 0; i < valueType.NumField(); i++ {
		tag := valueType.Field(i).Tag.Get("fixed")
		f := value.Field(i)

		if !f.CanInterface() || tag == "-" {
			continue
		}

		options := strings.Split(tag, ",")

		switch v := f.Interface().(type) {
		case bool:
			e.encodeBool(v, options)

		case time.Time:
			if err := e.encodeTime(v, options); err != nil {
				return nil, err
			}

		default:
			e.encodeValue(v, options)
		}
	}

	return e.Bytes(), nil
}

func (e *encoder) encodeBool(b bool, opts []string) error {
	var (
		chars  []string
		format string
	)

	switch len(opts) {
	case 3:
		format = opts[0]
		opts = opts[1:]
		fallthrough
	case 2:
		chars = opts
	default:
		return fmt.Errorf("invalid options")
	}

	i := 0

	if b {
		i = 1
	}

	_, err := e.WriteString(fmt.Sprintf(format, chars[i]))
	return err
}

func (e *encoder) encodeTime(t time.Time, opts []string) error {
	var (
		format string
		zero   bool
	)

	switch len(opts) {
	case 2:
		zero = opts[1] == "zero"
		fallthrough
	case 1:
		format = opts[0]
	default:
		return errInvalidDateFormat
	}

	if t.IsZero() && zero {
		e.WriteString(strings.Repeat("0", len(format)))
	} else {
		e.WriteString(t.Format(format))
	}

	return nil
}

func (e *encoder) encodeValue(v interface{}, opts []string) error {
	if len(opts) == 0 {
		return fmt.Errorf("invalid options")
	}

	_, err := e.WriteString(fmt.Sprintf(opts[0], v))
	return err
}

func Encode(v interface{}) ([]byte, error) {
	return new(encoder).encode(v)
}
