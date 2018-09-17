package rocket

type Response struct {
	headers []Header
	body    interface{}
}

type Header struct {
	Key   string
	Value string
}

func NewResponse(body interface{}) *Response {
	return &Response{
		headers: make([]Header, 0),
		body:    body,
	}
}

func (res *Response) Headers(headers ...Header) *Response {
	res.headers = append(res.headers, headers...)
	return res
}
