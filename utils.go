package hqtrivia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var apiURL string

func SetAPIUrl(Url string) {
	apiURL = Url
}

// Given a request pointer req, a boolean basBody (indicating whether the request has a body), and a token bearerToken (or "" for none),
// This function sets the headers of the provided request, to be in cooperation with the requirements for HQ Trivia
func setHeaders(req *http.Request, hasBody bool, bearerToken string) {
	req.Header.Add("x-hq-client", "iOS/1.8.0 b0")                         // add the client header (required)
	req.Header.Add("User-Agent", "HQ-iOS/0 CFNetwork/1390 Darwin/22.0.0") // add the user agent header (required)
	if hasBody {                                                          // add the content type header if needed
		req.Header.Add("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	}
}

// Given a response pointer res, this function deserializes the response body and returns a hashmap from strings to interfaces
func responseBody(res *http.Response) (map[string]interface{}, error) {
	body, err := io.ReadAll(res.Body) // try converting the response body into bytes
	if err != nil {
		return nil, err
	}
	var deserializedBody map[string]interface{}
	err = json.Unmarshal(body, &deserializedBody) // try converting the bytes to kv pairs
	if err != nil {
		return nil, err
	}
	return deserializedBody, nil
}

// Given an HTTP method, endpoint, some array of bytes body, and a token bearerToken (or "" for none),
// This function will execute the request and return the deserialized json in a hashmap from strings to interfaces, as well as the HTTP status code
func request(method string, endpoint string, body []byte, bearerToken string) (map[string]interface{}, int, error) {
	var req *http.Request
	var err error
	if len(body) > 1 {
		req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}
	if err != nil {
		return nil, -1, err
	}
	setHeaders(req, body != nil, bearerToken)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, -1, err
	}
	respBody, err := responseBody(response)
	if err != nil {
		return nil, -1, err
	}
	return respBody, response.StatusCode, nil
}
