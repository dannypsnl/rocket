package rocket

type Response struct {
	contentType string
	messages    []string
}

type ResponseBuilder struct {
	contentType string
	messages    []string
}

func (rb *ResponseBuilder) Done() Response {
	defer func() {
		rb.contentType = ""
		rb.messages = []string{}
	}()
	return Response{
		contentType: rb.contentType,
		messages:    rb.messages,
	}
}

func (rb *ResponseBuilder) ContentType(contentType string) *ResponseBuilder {
	rb.contentType = contentType
	return rb
}
