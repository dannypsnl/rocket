package rocket

type Response struct {
	headers map[string]string
	body    interface{}
}

type Headers map[string]string

func NewResponse(body interface{}) *Response {
	return &Response{
		headers: make(map[string]string),
		body:    body,
	}
}

func (res *Response) Headers(headers Headers) *Response {
	for k, v := range headers {
		res.headers[k] = v
	}
	return res
}
