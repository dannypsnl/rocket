package rocket

import (
	"testing"
)

func TestMatchArrayLen(t *testing.T) {
	ma := matchArray{"/", "/aabc", "/home/index", "/home/about"}
	for i := 0; i < 100; i++ {
		if len(ma) != ma.Len() {
			t.Error(`match array's Len is wrong`)
		}
	}
}

func TestMatchArrayLess(t *testing.T) {
	ma := matchArray{"/", "/home/index"}
	for i := 0; i < 100; i++ {
		if ma.Less(0, 1) {
			t.Error(`match array's Less is wrong`)
		}
	}
}
