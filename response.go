package rocket

type Response struct {
	headers map[string]string
	body    interface{}
}

type Header struct {
	Key   string
	Value string
}

func NewResponse(body interface{}) *Response {
	return &Response{
		headers: make(map[string]string),
		body:    body,
	}
}

func (res *Response) Headers(headers ...Header) *Response {
	for _, header := range headers {
		res.headers[header.Key] = header.Value
	}
	return res
}
