package rocket

type response struct {
	contentType string
	messages    []string
}

type ResponseBuilder struct {
	contentType string
	messages    []string
}

func (rb *ResponseBuilder) Done() response {
	defer func() {
		rb.contentType = ""
		rb.messages = []string{}
	}()
	return response{
		contentType: rb.contentType,
		messages:    rb.messages,
	}
}

func (rb *ResponseBuilder) ContentType(contentType string) *ResponseBuilder {
	rb.contentType = contentType
	return rb
}
