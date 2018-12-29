package rocket

import (
	"reflect"
	"testing"

	"github.com/dannypsnl/assert"
)

func TestParseParameter(t *testing.T) {
	assert := assert.NewTester(t)
	testCases := []struct {
		name          string
		paramStr      string
		expectedValue interface{}
	}{
		{
			"String",
			"change",
			"change",
		},
		{
			"Bool",
			"false",
			false,
		},
		{
			"Int",
			"50",
			50,
		},
		{
			"Int8",
			"50",
			int8(50),
		},
		{
			"Int16",
			"50",
			int16(50),
		},
		{
			"Int32",
			"50",
			int32(50),
		},
		{
			"Int64",
			"50",
			int64(50),
		},
		{
			"Uint",
			"3",
			uint(3),
		},
		{
			"Uint8",
			"3",
			uint8(3),
		},
		{
			"Uint16",
			"3",
			uint16(3),
		},
		{
			"Uint32",
			"3",
			uint32(3),
		},
		{
			"Uint64",
			"3",
			uint64(3),
		},
		{
			"Float32",
			"3",
			float32(3),
		},
		{
			"Float64",
			"3",
			float64(3),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			v, err := parseParameter(reflect.TypeOf(testCase.expectedValue), testCase.paramStr)
			assert.Eq(err, nil)
			assert.Eq(v.Interface(), testCase.expectedValue)
		})
	}
}
