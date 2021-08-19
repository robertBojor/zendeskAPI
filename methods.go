package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (z *API) Configure(baseURL, token string) {
	z.token = token
	z.baseURL = baseURL
	z.Error = nil
	z.Request = nil
	z.Response = nil
	z.ResponseBytes = make([]byte, 0)
}

func (z *API) NewRequest(method, endpoint string, payload interface{}) *API {
	prefixSlash := ""
	if !strings.HasPrefix(endpoint, "/") {
		prefixSlash = "/"
	}
	url := fmt.Sprintf("%s%s%s", z.baseURL, prefixSlash, endpoint)
	z.Error = nil
	var body *bytes.Buffer
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			z.Error = err
			return z
		}
		body = bytes.NewBuffer(payloadBytes)
	}
	if body == nil {
		z.Request, z.Error = http.NewRequest(method, url, nil)
	} else {
		z.Request, z.Error = http.NewRequest(method, url, body)
	}
	if z.Error != nil {
		return z
	}
	z.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.token))
	z.Request.Header.Add("Accept", "application/json")
	z.Request.Header.Add("Content-Type", "application/json")
	z.Request.Header.Add("Content-Language", "en")
	return z
}

func (z *API) Execute() {
	if z.Request == nil {
		z.Error = errors.New("no valid request found")
		return
	}
	client := &http.Client{}
	z.Response, z.Error = client.Do(z.Request)
	if z.Error != nil {
		z.Error = errors.Wrap(z.Error, "failed to execute request")
		return
	}
	z.ResponseBytes, z.Error = ioutil.ReadAll(z.Response.Body)
	if z.Error != nil {
		z.Error = errors.Wrap(z.Error, "failed to read response body")
		return
	}
	return
}
