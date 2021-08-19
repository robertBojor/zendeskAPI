package zendesk

import "net/http"

type API struct {
	token     string
	baseURL string
	Error         error
	Request       *http.Request
	Response      *http.Response
	ResponseBytes []byte
}