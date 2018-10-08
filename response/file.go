package response

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FileResponser struct {
	filepath string

	resp *Response
	err  error
}

func File(filepath string) *FileResponser {
	r := &FileResponser{
		filepath: filepath,
	}
	f, err := os.Open(filepath)
	if err != nil {
		r.err = err
		return r
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		r.err = err
		return r
	}
	r.resp = New(string(b))
	return r
}

func (r *FileResponser) ByFileSuffix(contentTypes map[string]string) *Response {
	if r.err != nil {
		r.resp.Status(http.StatusUnprocessableEntity)
	} else {
		r.resp.Status(http.StatusOK)
	}
	headers := Headers{}
	i := strings.LastIndexByte(r.filepath, '.')
	v, ok := contentTypes[r.filepath[i+1:]]
	if !ok {
		v = "text/plain"
	}
	headers["Content-Type"] = v

	r.resp.WithHeaders(headers)
	return r.resp
}

var DefaultContentTypes = map[string]string{
	"html": "text/html",
	"css":  "text/css",
	"js":   "application/javascript",
	"json": "application/json",
	"xml":  "application/xml",
	"gif":  "image/gif",
	"png":  "image/png",
	"jpeg": "image/jpeg",
}
