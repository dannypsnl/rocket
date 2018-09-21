package response

type Response struct {
	Headers map[string]string
	Body    interface{}
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		Headers: make(map[string]string),
		Body:    body,
	}
}

func (res *Response) WithHeaders(headers Headers) *Response {
	for k, v := range headers {
		res.Headers[k] = v
	}
	return res
}
