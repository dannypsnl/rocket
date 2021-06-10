package response

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
		return ContentTypeHTML
	case Json:
		return ContentTypeJSON
	default:
		return ContentTypeTextPlain
	}
}

const (
	ContentTypeHTML      = "text/html"
	ContentTypeJSON      = "application/json"
	ContentTypeTextPlain = "text/plain"
)
