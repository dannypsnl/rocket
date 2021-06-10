package response

import (
	"github.com/dannypsnl/rocket/response/content_type"
)

type (
	// Html is a mark for text is HTML
	// return "Content-Type": "text/html"
	Html string
	// Json will return "Content-Type": "application/json"
	Json string
)

func contentTypeOf(response interface{}) string {
	switch response.(type) {
	case Html:
		return content_type.HTML
	case Json:
		return content_type.JSON
	default:
		return content_type.TextPlain
	}
}
