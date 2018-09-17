package rocket

import (
	"net/http"
)

type Header struct {
	req *http.Request
}

func (h *Header) Get(key string) string {
	return h.req.Header.Get(key)
}
