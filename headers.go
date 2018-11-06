package rocket

import (
	"net/http"
)

type Headers struct {
	header http.Header
}

func (h *Headers) Get(key string) string {
	return h.header.Get(key)
}
