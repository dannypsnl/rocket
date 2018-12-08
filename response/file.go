package response

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type FileResponder struct {
	filepath string

	resp *Response
	err  error
}

func newFileResponder(filepath string) *FileResponder {
	return &FileResponder{
		filepath: filepath,
	}
}

func File(filepath string) *FileResponder {
	r := newFileResponder(filepath)
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

func (r *FileResponder) SetContentType(contentType ContentType) *Response {
	if r.err != nil {
		r.resp.Status(http.StatusUnprocessableEntity)
		return r.resp
	}
	r.resp.Status(http.StatusOK)
	headers := Headers{}
	v, ok := contentType.ByFileName(r.filepath)
	if !ok {
		v = "text/plain"
	}
	headers["Content-Type"] = v

	r.resp.Headers(headers)
	return r.resp
}

type ContentType interface {
	// ByFileName would receive a fileName, return contentType
	// if you don't know how to handle the fileName then just return `false`
	// it would fall back to `text/plain`
	ByFileName(fileName string) (contentType string, handled bool)
}

func ByFileNameSuffix() *fileSuffix {
	return &fileSuffix{}
}

type fileSuffix struct{}

func (f *fileSuffix) ByFileName(fileName string) (string, bool) {
	fileSuffix := fileName[strings.LastIndexByte(fileName, '.')+1:]
	contentType, handled := defaultSuffixMapToContentTypes[fileSuffix]
	return contentType, handled
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
