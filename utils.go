package hqtrivia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var apiURL string

func SetAPIUrl(Url string) {
	apiURL = Url
}

// Given an error, this function forces the error (if existant to cause a fatal error), with the option of a custom message str
func fatalError(err error, str string) {
	if err != nil {
		if str == "" {
			str = err.Error()
		}
		log.Fatalf("An error marked as fatal has occurred: %s", str)
	}
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
func responseBody(res *http.Response) map[string]interface{} {
	body, err := io.ReadAll(res.Body) // try converting the response body into bytes
	fatalError(err, "Byte Conversion (Response Body)")
	var deserializedBody map[string]interface{}
	err = json.Unmarshal(body, &deserializedBody) // try converting the bytes to kv pairs
	fatalError(err, "JSON Deserialization (Response Body)")
	return deserializedBody
}

// Given an HTTP method, endpoint, some array of bytes body, and a token bearerToken (or "" for none),
// This function will execute the request and return the deserialized json in a hashmap from strings to interfaces, as well as the HTTP status code
func request(method string, endpoint string, body []byte, bearerToken string) (map[string]interface{}, int) {
	var req *http.Request
	var err error
	if len(body) > 1 {
		req, err = http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, endpoint, nil)
	}
	fatalError(err, "HTTP Request (Creation)")
	setHeaders(req, body != nil, bearerToken)
	response, err := http.DefaultClient.Do(req)
	fatalError(err, "HTTP Request (Response)")
	return responseBody(response), response.StatusCode
}
