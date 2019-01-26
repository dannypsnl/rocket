package response

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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
	fileSuffix := filepath[strings.LastIndexByte(filepath, '.')+1:]
	return resp.ContentType(defaultSuffixMapToContentTypes[fileSuffix])
}

var defaultSuffixMapToContentTypes = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"gif":  "image/gif",
	"png":  "image/png",
	"jpeg": "image/jpeg",
}
