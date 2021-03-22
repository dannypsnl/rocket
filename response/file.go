package response

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

// File would create a response for file resource by provided filepath
func File(filepath string) *Response {
	var err error
	f, err := os.Open(filepath)
	if err != nil {
		return New("").Status(http.StatusNotFound)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return New("").Status(http.StatusUnprocessableEntity)
	}
	resp := New(string(b))
	return resp.ContentType(mimetype.Detect(b).String())
}
