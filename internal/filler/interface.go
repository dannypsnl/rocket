package filler

import (
	"reflect"
)

type Filler interface {
	Fill(ctx reflect.Value) error
}
