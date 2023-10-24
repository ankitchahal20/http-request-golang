package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ankit/project/http-request-golang/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	endpoint := fmt.Sprintf("%s/some-endpoint", "www.some-domain.com")
	ctx := &gin.Context{
		Request: &http.Request{
			Header: http.Header{
				constants.Authorization: []string{"Bearer some-random-beraer-token"},
			}},
	}
	token := ctx.Request.Header.Get(constants.Authorization) // token will be in middleware, here we are just fetching it from request header and using

	headers := map[string]string{
		constants.Authorization: token,
	}
	requestFields := ""

	body, marshallErr := json.Marshal(requestFields)
	if marshallErr != nil {
		log.Fatalf(fmt.Sprintf("error while marshaling request body to json : %v ", marshallErr), http.StatusInternalServerError)
	}

	paramHTTP := HTTPRequestParams{
		Headers: headers,
		URL:     endpoint,
		Body:    body,
		Method:  http.MethodPost,
	}
	transactionID := uuid.NewString()
	resp, err := paramHTTP.SendRequest(transactionID)
	if err != nil {
		log.Fatalf(fmt.Sprintf("error return from %v endpoint : %v", endpoint, err.Message), err.Code)
	} else if resp.StatusCode != http.StatusAccepted {
		log.Fatalf(fmt.Sprintf("received unexpected response status code from %v endpoint : %v", endpoint, err.Message), resp.StatusCode)
	}
	defer resp.Body.Close()
}
