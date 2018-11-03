package rocket

import (
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

func parseParameter(v reflect.Value, param string) reflect.Value {
	switch v.Kind() {
	case reflect.Bool:
		r, err := strconv.ParseBool(param)
		if err != nil {
			return reflect.ValueOf((*bool)(nil))
		}
		return reflect.ValueOf(r)
	case reflect.Int:
		r, err := strconv.ParseInt(param, 10, 0)
		if err != nil {
			return reflect.ValueOf((*int)(nil))
		}
		return reflect.ValueOf(int(r))
	case reflect.Int8:
		r, err := strconv.ParseInt(param, 10, 8)
		if err != nil {
			return reflect.ValueOf((*int8)(nil))
		}
		return reflect.ValueOf(int8(r))
	case reflect.Int16:
		r, err := strconv.ParseInt(param, 10, 16)
		if err != nil {
			return reflect.ValueOf((*int16)(nil))
		}
		return reflect.ValueOf(int16(r))
	case reflect.Int32:
		r, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			return reflect.ValueOf((*int32)(nil))
		}
		return reflect.ValueOf(int32(r))
	case reflect.Int64:
		r, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return reflect.ValueOf((*int64)(nil))
		}
		return reflect.ValueOf(int64(r))
	case reflect.Uint:
		r, err := strconv.ParseUint(param, 10, 0)
		if err != nil {
			return reflect.ValueOf((*uint)(nil))
		}
		return reflect.ValueOf(uint(r))
	case reflect.Uint8:
		r, err := strconv.ParseUint(param, 10, 8)
		if err != nil {
			return reflect.ValueOf((*uint8)(nil))
		}
		return reflect.ValueOf(uint8(r))
	case reflect.Uint16:
		r, err := strconv.ParseUint(param, 10, 16)
		if err != nil {
			return reflect.ValueOf((*uint16)(nil))
		}
		return reflect.ValueOf(uint16(r))
	case reflect.Uint32:
		r, err := strconv.ParseUint(param, 10, 32)
		if err != nil {
			return reflect.ValueOf((*uint32)(nil))
		}
		return reflect.ValueOf(uint32(r))
	case reflect.Uint64:
		r, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			return reflect.ValueOf((*uint64)(nil))
		}
		return reflect.ValueOf(r)
	case reflect.Float32:
		r, err := strconv.ParseFloat(param, 32)
		if err != nil {
			return reflect.ValueOf((*float32)(nil))
		}
		return reflect.ValueOf(float32(r))
	case reflect.Float64:
		r, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return reflect.ValueOf((*float64)(nil))
		}
		return reflect.ValueOf(r)
	case reflect.String:
		return reflect.ValueOf(param)
	default:
		panic("")
	}
}
