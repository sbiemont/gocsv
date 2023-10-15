package internal

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"

	"github.com/sbiemont/gocsv/lib"
)

// unmarshal a string into a reflect value
type unmarshaler func(string, reflect.Value) error

// first check that the string can be unmarshaled before processing
type unmarshalerWithCheck func(string, reflect.Value) (bool, error)

// convert to unmarshaler (without the check)
func (it unmarshalerWithCheck) toUnmarshaler() unmarshaler {
	return func(s string, v reflect.Value) error {
		_, err := it(s, v)
		return err
	}
}

// CacheUnmarshaler is used to store the unmarshaler function for the ith field
type CacheUnmarshaler map[int]unmarshaler

// NewCacheUnmarshaler init a new cache
func NewCacheUnmarshaler() CacheUnmarshaler {
	return make(CacheUnmarshaler)
}

func (c CacheUnmarshaler) storeWithCheck(col int, s string, v reflect.Value, fct unmarshalerWithCheck) (bool, error) {
	// Use fct and store it in cache if check is ok
	ok, err := fct(s, v)
	switch {
	case err != nil:
		return false, err
	case ok:
		c[col] = fct.toUnmarshaler()
		return true, nil
	default:
		return false, nil
	}
}

func (c CacheUnmarshaler) storeWithChecks(
	col int,
	s string,
	v reflect.Value,
	fcts []unmarshalerWithCheck,
) (bool, error) {
	for _, fct := range fcts {
		ok, err := c.storeWithCheck(col, s, v, fct)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func (c CacheUnmarshaler) store(col int, s string, v reflect.Value, fct unmarshaler) error {
	// Use fct and store it in cache if check is ok
	err := fct(s, v)
	if err != nil {
		return err
	}
	c[col] = fct
	return nil
}

// Use cache if present
func (c CacheUnmarshaler) use(col int, s string, v reflect.Value) (bool, error) {
	fctCache, ok := c[col]
	if ok {
		return true, fctCache(s, v)
	}
	return false, nil
}

// Unmarhsal a list of fields in the given instance
func Unmarshal[T any](ct CacheTags[T], cm CacheUnmarshaler, inputs []string, item *T) error {
	// Error helper using item column
	makeErr := func(col int, e error) error {
		return fmt.Errorf("col %d: %w", col, e)
	}

	val := reflect.Indirect(reflect.ValueOf(item))
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		// Get "csv" info => parse `csv:"tag0,tag1,..,tagN"`
		tag, ok := ct.tag(i)
		if ok {
			col := tag.col
			omitEmpty := tag.omitEmpty

			if col >= len(inputs) {
				if omitEmpty {
					continue
				}
				return fmt.Errorf("column %d out of bounds", col)
			}

			// Fetch data
			input := inputs[col]
			if input == "" && omitEmpty { // omit empty
				continue
			}

			// Fetch current attribute
			field := val.Field(i)

			// If pointer, allocate a new object and use it
			if field.Type().Kind() == reflect.Ptr {
				field.Set(reflect.New(field.Type().Elem()))
				field = field.Elem()
			}

			// Use cache if filled
			ok, err := cm.use(col, input, field)
			switch {
			case err != nil:
				return makeErr(col, err)
			case ok:
				continue
			}

			// Unmarshalers with checks
			ok, err = cm.storeWithChecks(col, input, field, unmarshalersWithCheckConfig)
			switch {
			case err != nil:
				return makeErr(col, err)
			case ok:
				continue
			}

			// Choose the field type
			k := field.Type().Kind()
			unmarshal, ok := unmarshalersConfig[k]
			if !ok {
				return makeErr(col, fmt.Errorf("unknown type %s", k))
			}
			err = cm.store(col, input, field, unmarshal)
			if err != nil {
				return makeErr(col, err)
			}
		}
	}
	return nil
}

var unmarshalersWithCheckConfig = []unmarshalerWithCheck{
	csvUnmarshaler,
	textUnmarshaler,
}

var unmarshalersConfig = map[reflect.Kind]unmarshaler{
	reflect.Int:   intUnmarshaler,
	reflect.Int8:  intUnmarshaler,
	reflect.Int16: intUnmarshaler,
	reflect.Int32: intUnmarshaler,
	reflect.Int64: intUnmarshaler,

	reflect.Uint:   uintUnmarshaler,
	reflect.Uint8:  uintUnmarshaler,
	reflect.Uint16: uintUnmarshaler,
	reflect.Uint32: uintUnmarshaler,
	reflect.Uint64: uintUnmarshaler,

	reflect.Float32: floatUnmarshaler,
	reflect.Float64: floatUnmarshaler,

	reflect.String: stringUnmarshaler,

	reflect.Bool: boolUnmarshaler,
}

// Check for csv unmarshaler
func csvUnmarshaler(in string, field reflect.Value) (bool, error) {
	u, ok := field.Addr().Interface().(lib.Unmarshaler)
	if !ok {
		return false, nil
	}
	return true, u.UnmarshalCSV(in)
}

// Check for text unmarshaler
func textUnmarshaler(in string, field reflect.Value) (bool, error) {
	u, ok := field.Addr().Interface().(encoding.TextUnmarshaler)
	if !ok {
		return false, nil
	}

	return true, u.UnmarshalText([]byte(in))
}

func intUnmarshaler(in string, field reflect.Value) error {
	res, err := strconv.ParseInt(in, 0, 64)
	if err != nil {
		return err
	}
	field.SetInt(res)
	return nil
}

func uintUnmarshaler(in string, field reflect.Value) error {
	res, err := strconv.ParseUint(in, 0, 64)
	if err != nil {
		return err
	}
	field.SetUint(res)
	return nil
}

func floatUnmarshaler(in string, field reflect.Value) error {
	res, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return err
	}
	field.SetFloat(res)
	return nil
}

func stringUnmarshaler(in string, field reflect.Value) error {
	field.SetString(in)
	return nil
}

func boolUnmarshaler(in string, field reflect.Value) error {
	field.SetBool(in == "true")
	return nil
}
