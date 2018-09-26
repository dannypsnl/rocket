package response

func ContentTypeOf(response interface{}) string {
	switch response.(type) {
	case Html:
		return "text/html"
	case Json:
		return "application/json"
	case string:
		return "text/plain"
	default:
		return "text/plain"
	}
}

// Html is a mark for text is HTML
// return "Content-Type": "text/html"
type Html string

// Json will return "Content-Type": "application/json"
type Json string
