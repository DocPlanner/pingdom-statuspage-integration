package mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Request struct {
	Method  string
	Host    string
	URLPath string
}

type Response struct {
	Body        string
	Code        int
	ContentType string
}

type transport struct {
	http.RoundTripper
	Responses map[Request]Response
}

func NewTransport(responses *map[Request]Response) *transport {
	return &transport{
		Responses: *responses,
	}
}

func (t *transport) RoundTrip(req *http.Request) (response *http.Response, err error) {
	var currentResponse *Response

	currentResponse = t.FindFirstResponseForMatchingRequest(req)
	if currentResponse == nil {
		return nil, fmt.Errorf("Response for following request not found: %s %s %s", req.Method, req.Host, req.URL.Path)
	}

	if currentResponse.Code == 0 {
		currentResponse.Code = http.StatusOK
	}

	if currentResponse.Body == "" {
		currentResponse.Body = "{}"
	}

	if currentResponse.ContentType == "" {
		currentResponse.ContentType = "application/json"
	}

	response = &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: currentResponse.Code,
	}
	response.Header.Set("Content-Type", currentResponse.ContentType)

	response.Body = ioutil.NopCloser(strings.NewReader(currentResponse.Body))

	return response, err
}

func (t *transport) FindFirstResponseForMatchingRequest(request *http.Request) *Response {
	for req, rsp := range t.Responses {
		if request.Method == req.Method && request.Host == req.Host && request.URL.Path == req.URLPath {
			delete(t.Responses, req)
			return &rsp
		}
	}

	return nil
}
