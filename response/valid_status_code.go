// +build go1.11

package response

import (
	"net/http"
)

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
