package rocket

import (
	"net/http"
)

type Headers struct {
	req *http.Request
}

func (h *Headers) Get(key string) string {
	return h.req.Header.Get(key)
}
