package internal

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/sbiemont/gocsv/lib"
)

// marshal a reflect value into a string
type marshaler func(reflect.Value) (string, error)

// first check that the value can be marshaled before processing
type marshalerWithCheck func(reflect.Value) (string, bool, error)

// convert to unmarshaler (without the check)
func (it marshalerWithCheck) toMarshaler() marshaler {
	return func(v reflect.Value) (string, error) {
		s, _, err := it(v)
		return s, err
	}
}

// CacheMarshaler store a marshaler for the ith field
type CacheMarshaler map[int]marshaler

// NewCacheMarshaler init a new cache
func NewCacheMarshaler() CacheMarshaler {
	return make(CacheMarshaler)
}

func (c CacheMarshaler) use(col int, v reflect.Value) (string, bool, error) {
	fctCache, ok := c[col]
	if ok {
		s, err := fctCache(v)
		return s, true, err
	}
	return "", false, nil
}

func (c CacheMarshaler) storeWithCheck(col int, v reflect.Value, fct marshalerWithCheck) (string, bool, error) {
	// Use fct and store it in cache if check is ok
	s, ok, err := fct(v)
	switch {
	case err != nil:
		return "", false, err
	case ok:
		c[col] = fct.toMarshaler()
		return s, true, nil
	default:
		return "", false, nil
	}
}

func (c CacheMarshaler) storeWithChecks(col int, v reflect.Value, fcts []marshalerWithCheck) (string, bool, error) {
	for _, fct := range fcts {
		s, ok, err := c.storeWithCheck(col, v, fct)
		if ok || err != nil {
			return s, ok, err
		}
	}
	return "", false, nil
}

func (c CacheMarshaler) store(col int, v reflect.Value, fct marshaler) (string, error) {
	// Use fct and store it in cache if check is ok
	s, err := fct(v)
	if err != nil {
		return "", err
	}
	c[col] = fct
	return s, nil
}

// Marshal a given struct into a list of csv values
func Marshal[T any](ct CacheTags[T], cm CacheMarshaler, item T) ([]string, error) {
	// Error helper using item column
	makeErr := func(col int, e error) ([]string, error) {
		return nil, fmt.Errorf("col %d: %w", col, e)
	}

	outputs := make(map[int]string)
	val := reflect.Indirect(reflect.ValueOf(item))
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Get "csv" info => parse `csv:"tag0,tag1,..,tagN"`
		tag, ok := ct.tag(i)
		if ok {
			// Fetch tags (the first one is the column)
			col := tag.col
			omitEmpty := tag.omitEmpty

			// Fetch current attribute
			field := val.Field(i)

			// If pointer, controls is nil and omit empty
			if field.Type().Kind() == reflect.Ptr {
				isNil := field.IsNil()
				switch {
				case isNil && omitEmpty:
					outputs[col] = ""
					continue
				case isNil && !omitEmpty:
					return makeErr(col, fmt.Errorf("nil value found"))
				default: // field.IsNil(): false
					field = field.Elem()
				}
			}

			// Use cache if filled
			res, ok, err := cm.use(col, field)
			switch {
			case err != nil:
				return makeErr(col, err)
			case ok:
				outputs[col] = res
				continue
			}

			// Ordered marshalers with check
			res, ok, err = cm.storeWithChecks(col, field, marshalersWithCheckConfig)
			switch {
			case err != nil:
				return makeErr(col, err)
			case ok:
				outputs[col] = res
				continue
			}

			// Choose the field type
			k := field.Type().Kind()
			marshal, ok := marshalersConfig[k]
			if !ok {
				return makeErr(col, fmt.Errorf("unknown type %s", k))
			}
			res, err = cm.store(col, field, marshal)
			if err != nil {
				return makeErr(col, err)
			}
			outputs[col] = res
		}
	}

	// Find max col
	maxCol := ct.maxCol()

	// Convert output to result
	result := make([]string, maxCol+1)
	for idx, output := range outputs {
		result[idx] = output
	}
	return result, nil
}

var marshalersWithCheckConfig = []marshalerWithCheck{
	csvMarshaler,
	textMarshaler,
}

var marshalersConfig = map[reflect.Kind]marshaler{
	reflect.Int:   intMarshaler,
	reflect.Int8:  intMarshaler,
	reflect.Int16: intMarshaler,
	reflect.Int32: intMarshaler,
	reflect.Int64: intMarshaler,

	reflect.Uint:   uintMarshaler,
	reflect.Uint8:  uintMarshaler,
	reflect.Uint16: uintMarshaler,
	reflect.Uint32: uintMarshaler,
	reflect.Uint64: uintMarshaler,

	reflect.Float32: floatMarshaler,
	reflect.Float64: floatMarshaler,

	reflect.String: stringMarshaler,

	reflect.Bool: boolMarshaler,
}

// Check for csv marshaler
func csvMarshaler(field reflect.Value) (string, bool, error) {
	// u, ok := field.Addr().Interface().(Marshaler)
	u, ok := field.Interface().(lib.Marshaler)
	if !ok {
		return "", false, nil
	}
	res, err := u.MarshalCSV()
	return res, true, err
}

// Check for text marshaler
func textMarshaler(field reflect.Value) (string, bool, error) {
	// u, ok := field.Addr().Interface().(encoding.TextMarshaler)
	u, ok := field.Interface().(encoding.TextMarshaler)
	if !ok {
		return "", false, nil
	}
	res, err := u.MarshalText()
	return string(res), true, err
}

func intMarshaler(field reflect.Value) (string, error) {
	return fmt.Sprintf("%d", field.Int()), nil
}

func uintMarshaler(field reflect.Value) (string, error) {
	return fmt.Sprintf("%d", field.Uint()), nil
}

func floatMarshaler(field reflect.Value) (string, error) {
	return fmt.Sprintf("%f", field.Float()), nil
}

func stringMarshaler(field reflect.Value) (string, error) {
	return field.String(), nil
}

func boolMarshaler(field reflect.Value) (string, error) {
	if field.Bool() {
		return "true", nil
	}
	return "false", nil
}
