package rocket

import (
	"errors"
	"reflect"
	"strconv"
)

var (
	invalidType = errors.New("invalid type")
)

func parseParameter(vt reflect.Type, param string) (reflect.Value, error) {
	switch vt.Kind() {
	case reflect.Bool:
		r, _ := strconv.ParseBool(param)
		return reflect.ValueOf(r), nil
	case reflect.Int:
		r, _ := strconv.ParseInt(param, 10, 0)
		return reflect.ValueOf(int(r)), nil
	case reflect.Int8:
		r, _ := strconv.ParseInt(param, 10, 8)
		return reflect.ValueOf(int8(r)), nil
	case reflect.Int16:
		r, _ := strconv.ParseInt(param, 10, 16)
		return reflect.ValueOf(int16(r)), nil
	case reflect.Int32:
		r, _ := strconv.ParseInt(param, 10, 32)
		return reflect.ValueOf(int32(r)), nil
	case reflect.Int64:
		r, _ := strconv.ParseInt(param, 10, 64)
		return reflect.ValueOf(int64(r)), nil
	case reflect.Uint:
		r, _ := strconv.ParseUint(param, 10, 0)
		return reflect.ValueOf(uint(r)), nil
	case reflect.Uint8:
		r, _ := strconv.ParseUint(param, 10, 8)
		return reflect.ValueOf(uint8(r)), nil
	case reflect.Uint16:
		r, _ := strconv.ParseUint(param, 10, 16)
		return reflect.ValueOf(uint16(r)), nil
	case reflect.Uint32:
		r, _ := strconv.ParseUint(param, 10, 32)
		return reflect.ValueOf(uint32(r)), nil
	case reflect.Uint64:
		r, _ := strconv.ParseUint(param, 10, 64)
		return reflect.ValueOf(r), nil
	case reflect.Float32:
		r, _ := strconv.ParseFloat(param, 32)
		return reflect.ValueOf(float32(r)), nil
	case reflect.Float64:
		r, _ := strconv.ParseFloat(param, 64)
		return reflect.ValueOf(r), nil
	case reflect.String:
		return reflect.ValueOf(param), nil
	case reflect.Ptr:
		// We use pointer represents optional field
		parsedValue, err := parseParameter(vt.Elem(), param)
		if err != nil {
			return reflect.Zero(vt), err
		}
		ptrToVal := reflect.New(vt.Elem()) // ptrToVal := new(TypeOf(parsedValue))
		ptrToVal.Elem().Set(parsedValue)   // *ptrToVal = parsedValue
		return ptrToVal, nil
	default:
		panic("unsupported parameter type")
	}
}
