package file

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dannypsnl/rocket/response"
)

type Responser struct {
	resp *response.Response

	fileName string
	err      error
}

func Response(fileName string) *Responser {
	r := &Responser{
		fileName: fileName,
	}
	f, err := os.Open(fileName)
	if err != nil {
		r.err = err
		return r
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		r.err = err
		return r
	}
	r.resp = response.New(string(b))
	return r
}

func (r *Responser) ByFileSuffix() *response.Response {
	if r.err != nil {
		r.resp.Status(http.StatusUnprocessableEntity)
	} else {
		r.resp.Status(http.StatusOK)
	}
	headers := response.Headers{}
	i := strings.LastIndexByte(r.fileName, '.')
	headers["Content-Type"] = contentTypeByFileSuffix(r.fileName[i+1:])

	r.resp.WithHeaders(headers)
	return r.resp
}

func contentTypeByFileSuffix(fileSuffix string) string {
	switch fileSuffix {
	case "html":
		return "text/html"
	case "css":
		return "text/css"
	case "js":
		return "application/javascript"
	case "json":
		return "application/json"
	case "xml":
		return "application/xml"
	case "gif":
		return "image/gif"
	case "png":
		return "image/png"
	case "jpeg":
		return "image/jpeg"
	default:
		return "text/plain"
	}
}
