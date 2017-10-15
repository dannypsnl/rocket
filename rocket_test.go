package rocket

import "testing"

func TestSplit(t *testing.T) {
	match, params := split("/")
	if match != "/" {
		t.Error(`match should be '/'.`, match)
	}
	if len(params) != 0 {
		t.Error(`params should be empty.`, params)
	}
}
