package rocket

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
)

func verifyBase(route string) bool {
	r, _ := regexp.Compile(".*?[:*].*?")
	// Contains : part will Match, it can be on a Base Route
	if r.MatchString(route) {
		panic("Base route can not contain dynamic route.")
	}
	return true
}

var (
	invaildType = errors.New("invalid type")
)

func parseParameter(vt reflect.Type, param string) (reflect.Value, error) {
	switch vt.Kind() {
	case reflect.Bool:
		r, err := strconv.ParseBool(param)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(r), nil
	case reflect.Int:
		r, err := strconv.ParseInt(param, 10, 0)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(int(r)), nil
	case reflect.Int8:
		r, err := strconv.ParseInt(param, 10, 8)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(int8(r)), nil
	case reflect.Int16:
		r, err := strconv.ParseInt(param, 10, 16)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(int16(r)), nil
	case reflect.Int32:
		r, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(int32(r)), nil
	case reflect.Int64:
		r, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(int64(r)), nil
	case reflect.Uint:
		r, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(uint(r)), nil
	case reflect.Uint8:
		r, err := strconv.ParseUint(param, 10, 8)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(uint8(r)), nil
	case reflect.Uint16:
		r, err := strconv.ParseUint(param, 10, 16)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(uint16(r)), nil
	case reflect.Uint32:
		r, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(uint32(r)), nil
	case reflect.Uint64:
		r, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(r), nil
	case reflect.Float32:
		r, err := strconv.ParseFloat(param, 32)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(float32(r)), nil
	case reflect.Float64:
		r, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return reflect.Zero(vt), invaildType
		}
		return reflect.ValueOf(r), nil
	case reflect.String:
		return reflect.ValueOf(param), nil
	case reflect.Ptr:
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
