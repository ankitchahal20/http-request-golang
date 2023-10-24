package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ankit/project/http-request-golang/constants"
)

type HTTPRequestParams struct {
	URL     string            `json:"url"`
	APIKey  string            `json:"apiKey,omitempty"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Token   string            `json:"token"`
	Body    []byte            `json:"body,omitempty"`
}

type Error struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Trace   string `json:"trace"`
}

type IHTTPRequest interface {
	SendRequest(*HTTPRequestParams, string) (http.Response, *Error)
}

func (httpRequestParam *HTTPRequestParams) SendRequest(transactionID string) (*http.Response, *Error) {
	var req *http.Request
	var err error
	// create http new request
	if httpRequestParam.Method == http.MethodGet {
		req, err = http.NewRequest(httpRequestParam.Method, httpRequestParam.URL, http.NoBody)
	} else if httpRequestParam.Method == http.MethodPost {
		req, err = http.NewRequest(httpRequestParam.Method, httpRequestParam.URL, bytes.NewBuffer(httpRequestParam.Body))
	}

	if err != nil {

		return &http.Response{}, &Error{
			Status:  http.StatusServiceUnavailable,
			Message: fmt.Sprintf("could not create http request : %v", err),
		}
	}

	req.Header.Add(constants.ContentType, constants.JSONApplication)

	if httpRequestParam.Token != "" {
		req.Header.Add(constants.Authorization, "Bearer "+httpRequestParam.Token)
	}

	// appending vital headers
	for key, value := range httpRequestParam.Headers {
		req.Header.Set(key, value)
	}
	// append key when required
	if httpRequestParam.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+httpRequestParam.APIKey)
	}

	// get response from the api
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, &Error{
			Status:  resp.StatusCode,
			Message: fmt.Sprintf("failed to get the http response : %v", err.Error()),
		}
	}
	return resp, nil
}
