package zendesk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

func (z *API) Configure(baseURL, token string) {
	z.token = token
	z.debug = false
	z.baseURL = baseURL
	z.Error = nil
	z.Request = nil
	z.Response = nil
	z.ResponseBytes = make([]byte, 0)
}

func (z *API) SetDebug(debugState bool) {
	z.debug = debugState
}

func (z *API) createRequest(method, endpoint string, payload interface{}) *API {
	prefixSlash := ""
	if !strings.HasPrefix(endpoint, "/") {
		prefixSlash = "/"
	}
	url := fmt.Sprintf("%s%s%s", z.baseURL, prefixSlash, endpoint)
	if z.debug {
		logrus.Infof("✅ Created URL for new request: %s", url)
	}
	z.Error = nil
	var body *bytes.Buffer
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			z.Error = err
			return z
		}
		body = bytes.NewBuffer(payloadBytes)
		if z.debug {
			logrus.Infof("✅ Payload created: %s", string(payloadBytes))
		}
	}
	if body == nil {
		z.Request, z.Error = http.NewRequest(method, url, nil)
	} else {
		z.Request, z.Error = http.NewRequest(method, url, body)
	}
	if z.debug {
		logrus.Infof("✅ Created new request")
	}
	if z.Error != nil {
		if z.debug {
			logrus.Errorf("❌ Failed to create new request - %s", z.Error.Error())
		}
		return z
	}
	z.Request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", z.token))
	z.Request.Header.Add("Accept", "application/json")
	z.Request.Header.Add("Content-Type", "application/json")
	z.Request.Header.Add("Content-Language", "en")
	return z
}

func (z *API) execute() {
	if z.Error != nil {
		return
	}
	if z.Request == nil {
		z.Error = errors.New("no request found")
		return
	}
	z.Error = nil
	client := &http.Client{}
	z.Response, z.Error = client.Do(z.Request)
	if z.debug {
		logrus.Infof("✅ Executing request...")
	}
	if z.Error != nil {
		if z.debug {
			logrus.Errorf("❌ Failed to execute request - %s", z.Error.Error())
		}
		z.Error = errors.Wrap(z.Error, "failed to execute request")
		return
	}
	z.ResponseBytes, z.Error = ioutil.ReadAll(z.Response.Body)
	if z.debug {
		logrus.Infof("✅ Response code %d - %s", z.Response.StatusCode, string(z.ResponseBytes))
	}
	if z.Error != nil {
		if z.debug {
			logrus.Errorf("❌ Failed to read response body - %s", z.Error.Error())
		}
		z.Error = errors.Wrap(z.Error, "failed to read response body")
		return
	}
	return
}
