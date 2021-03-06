package parse

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseParameter(t *testing.T) {
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
			v, err := ParseParameter(reflect.TypeOf(testCase.expectedValue), testCase.paramStr)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedValue, v.Interface())
		})
	}
}
