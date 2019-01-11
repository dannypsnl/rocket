package response

import (
	"fmt"
	"net/http"

	"github.com/dannypsnl/rocket/cookie"
)

type Response struct {
	headers map[string]string
	Body    interface{}

	cookies    []*http.Cookie
	statusCode int
}

type Headers map[string]string

func New(body interface{}) *Response {
	return &Response{
		headers: map[string]string{
			"Content-Type": contentTypeOf(body),
		},
		Body:    body,
		cookies: make([]*http.Cookie, 0),
	}
}

var validStatusCode = map[int]bool{
	http.StatusContinue:           true,
	http.StatusSwitchingProtocols: true,
	http.StatusProcessing:         true,

	http.StatusOK:                   true,
	http.StatusCreated:              true,
	http.StatusAccepted:             true,
	http.StatusNonAuthoritativeInfo: true,
	http.StatusNoContent:            true,
	http.StatusResetContent:         true,
	http.StatusPartialContent:       true,
	http.StatusMultiStatus:          true,
	http.StatusAlreadyReported:      true,
	http.StatusIMUsed:               true,

	http.StatusMultipleChoices:  true,
	http.StatusMovedPermanently: true,
	http.StatusFound:            true,
	http.StatusSeeOther:         true,
	http.StatusNotModified:      true,
	http.StatusUseProxy:         true,

	http.StatusTemporaryRedirect: true,
	http.StatusPermanentRedirect: true,

	http.StatusBadRequest:                   true,
	http.StatusUnauthorized:                 true,
	http.StatusPaymentRequired:              true,
	http.StatusForbidden:                    true,
	http.StatusNotFound:                     true,
	http.StatusMethodNotAllowed:             true,
	http.StatusNotAcceptable:                true,
	http.StatusProxyAuthRequired:            true,
	http.StatusRequestTimeout:               true,
	http.StatusConflict:                     true,
	http.StatusGone:                         true,
	http.StatusLengthRequired:               true,
	http.StatusPreconditionFailed:           true,
	http.StatusRequestEntityTooLarge:        true,
	http.StatusRequestURITooLong:            true,
	http.StatusUnsupportedMediaType:         true,
	http.StatusRequestedRangeNotSatisfiable: true,
	http.StatusExpectationFailed:            true,
	http.StatusTeapot:                       true,
	http.StatusMisdirectedRequest:           true,
	http.StatusUnprocessableEntity:          true,
	http.StatusLocked:                       true,
	http.StatusFailedDependency:             true,
	http.StatusUpgradeRequired:              true,
	http.StatusPreconditionRequired:         true,
	http.StatusTooManyRequests:              true,
	http.StatusRequestHeaderFieldsTooLarge:  true,
	http.StatusUnavailableForLegalReasons:   true,

	http.StatusInternalServerError:           true,
	http.StatusNotImplemented:                true,
	http.StatusBadGateway:                    true,
	http.StatusServiceUnavailable:            true,
	http.StatusGatewayTimeout:                true,
	http.StatusHTTPVersionNotSupported:       true,
	http.StatusVariantAlsoNegotiates:         true,
	http.StatusInsufficientStorage:           true,
	http.StatusLoopDetected:                  true,
	http.StatusNotExtended:                   true,
	http.StatusNetworkAuthenticationRequired: true,
}

func (res *Response) Status(code int) *Response {
	if !validStatusCode[code] {
		panic(fmt.Errorf("reject invalid status code: %d", code))
	}
	if res.statusCode != 0 {
		panic("status code already be changed, you can't modify it twice!")
	}
	res.statusCode = code
	return res
}

func (res *Response) Headers(headers Headers) *Response {
	for k, v := range headers {
		res.headers[k] = v
	}
	return res
}

func (res *Response) Cookies(cs ...*cookie.Cookie) *Response {
	for _, c := range cs {
		res.cookies = append(res.cookies, c.Generate())
	}
	return res
}

func (res *Response) WriteTo(w http.ResponseWriter) {
	res.setHeaders(w)
	res.setCookie(w)
	res.setStatusCode(w)
	fmt.Fprint(w, res.Body)
}

func (res *Response) setHeaders(w http.ResponseWriter) {
	for k, v := range res.headers {
		w.Header().Set(k, v)
	}
}

func (res *Response) setCookie(w http.ResponseWriter) {
	for _, c := range res.cookies {
		http.SetCookie(w, c)
	}
}

func (res *Response) setStatusCode(w http.ResponseWriter) {
	if res.statusCode != 0 {
		w.WriteHeader(res.statusCode)
	}
}
