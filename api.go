package zendesk

import "net/http"

type API struct {
	token         string
	baseURL       string
	debug         bool
	Error         error
	Request       *http.Request
	Response      *http.Response
	ResponseBytes []byte
}
