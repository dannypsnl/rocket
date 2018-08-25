package rocket

import (
	"github.com/dannypsnl/assert"
	"testing"

	"reflect"
)

func TestParseParameter(t *testing.T) {
	assert := assert.NewTester(t)
	t.Run("String", func(t *testing.T) {
		rv := reflect.ValueOf("string")
		nv := parseParameter(rv, "change")
		assert.Eq(nv.Interface(), "change")
	})
	t.Run("Bool", func(t *testing.T) {
		rv := reflect.ValueOf(true)
		nv := parseParameter(rv, "false")
		assert.Eq(nv.Interface(), false)
	})
	t.Run("Int", func(t *testing.T) {
		rv := reflect.ValueOf(10)
		nv := parseParameter(rv, "50")
		assert.Eq(nv.Interface(), 50)
		assert.Eq(nv.Kind(), reflect.Int)
	})
	t.Run("Int8", func(t *testing.T) {
		rv := reflect.ValueOf(int8(10))
		nv := parseParameter(rv, "50")
		assert.Eq(nv.Interface(), int8(50))
		assert.Eq(nv.Kind(), reflect.Int8)
	})
	t.Run("Int16", func(t *testing.T) {
		rv := reflect.ValueOf(int16(10))
		nv := parseParameter(rv, "50")
		assert.Eq(nv.Interface(), int16(50))
		assert.Eq(nv.Kind(), reflect.Int16)
	})
	t.Run("Int32", func(t *testing.T) {
		rv := reflect.ValueOf(int32(10))
		nv := parseParameter(rv, "50")
		assert.Eq(nv.Interface(), int32(50))
		assert.Eq(nv.Kind(), reflect.Int32)
	})
	t.Run("Int64", func(t *testing.T) {
		rv := reflect.ValueOf(int64(10))
		nv := parseParameter(rv, "50")
		assert.Eq(nv.Interface(), int64(50))
		assert.Eq(nv.Kind(), reflect.Int64)
	})
	t.Run("Uint", func(t *testing.T) {
		rv := reflect.ValueOf(uint(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), uint(3))
		assert.Eq(nv.Kind(), reflect.Uint)
	})
	t.Run("Uint8", func(t *testing.T) {
		rv := reflect.ValueOf(uint8(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), uint8(3))
		assert.Eq(nv.Kind(), reflect.Uint8)
	})
	t.Run("Uint16", func(t *testing.T) {
		rv := reflect.ValueOf(uint16(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), uint16(3))
		assert.Eq(nv.Kind(), reflect.Uint16)
	})
	t.Run("Uint32", func(t *testing.T) {
		rv := reflect.ValueOf(uint32(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), uint32(3))
		assert.Eq(nv.Kind(), reflect.Uint32)
	})
	t.Run("Uint64", func(t *testing.T) {
		rv := reflect.ValueOf(uint64(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), uint64(3))
		assert.Eq(nv.Kind(), reflect.Uint64)
	})
	t.Run("Float32", func(t *testing.T) {
		rv := reflect.ValueOf(float32(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), float32(3))
		assert.Eq(nv.Kind(), reflect.Float32)
	})
	t.Run("Float64", func(t *testing.T) {
		rv := reflect.ValueOf(float64(5))
		nv := parseParameter(rv, "3")
		assert.Eq(nv.Interface(), float64(3))
		assert.Eq(nv.Kind(), reflect.Float64)
	})
}
