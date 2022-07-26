package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

func CreateStringPointer(s string) *string {
	return &s
}

type Response struct {
	ApiVersion string           `json:"api_version,omitempty"`
	Kind       string           `json:"kind,omitempty"`
	Metadata   ResponseMetadata `json:"metadata,omitempty"`
	Data       interface{}      `json:"data,omitempty"`
}

type ResponseMetadata struct {
	First     string `json:"first,omitempty"`
	Last      string `json:"last,omitempty"`
	Prev      string `json:"prev,omitempty"`
	Next      string `json:"next,omitempty"`
	TotalSize int    `json:"totalSize,omitempty"`
}

type DoHttpRequestParameters struct {
	HttpClient       *http.Client
	Req              *http.Request
	ParameterSession Session
	DefaultSession   Session
}

type ErrorResponseEntity struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (e *ErrorResponseEntity) Error() string {
	return e.Message
}

type ErrorHandler interface {
	Handle(e error, response *http.Response) (*ErrorResponseEntity, error)
}

// type DoHTTPRequestResponse struct{}
func DoHttpRequest(req DoHttpRequestParameters, errorHandler ErrorHandler) (*http.Response, *ErrorResponseEntity, error) {
	err := SetAuthHeader(req.Req, req.ParameterSession, req.DefaultSession)
	if err != nil {
		return nil, nil, err
	}

	resp, err := req.HttpClient.Do(req.Req)
	if err != nil {
		return resp, nil, err
	}

	errorResponseEntity, err := errorHandler.Handle(err, resp)
	if err != nil {
		return nil, nil, err
	}

	return resp, errorResponseEntity, nil
}

func DeserializeResponse[T any](body io.ReadCloser) (T, error) {
	var payload T
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return payload, err
	}
	err = json.Unmarshal(bodyBytes, &payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

// func HandlePagedResponse[T any](r DoHttpRequestParameters, initialResp Response) ([]T, error) { // TODO: Need testing with statuscode
// 	payload := []T{}
// 	container := []interface{}{}
// 	container = append(container, initialResp.Data)

// 	nextUrl := initialResp.Metadata.Next

// 	for {
// 		if nextUrl == "" {
// 			break
// 		}
// 		fmt.Println(nextUrl)
// 		req, err := http.NewRequest("GET", nextUrl, nil)
// 		if err != nil {
// 			return payload, err
// 		}

// 		nextResp, err := DoHttpRequest(DoHttpRequestParameters{
// 			HttpClient:       r.HttpClient,
// 			Req:              req,
// 			ParameterSession: r.ParameterSession,
// 			DefaultSession:   r.DefaultSession,
// 		})
// 		if err != nil {
// 			return payload, err
// 		}

// 		// TODO: Deserialize

// 		container = append(container, nextResp.Data)
// 		nextUrl = nextResp.Metadata.Next
// 	}

// 	// Deserialise
// 	// TODO: Perhaps replace with https://github.com/mitchellh/mapstructure
// 	denested := []map[string]interface{}{}
// 	for _, c := range container {
// 		objects := c.([]interface{})
// 		for _, object := range objects {
// 			denested = append(denested, object.(map[string]interface{}))
// 		}
// 	}

// 	for _, d := range denested {
// 		var deserialised T
// 		serialised, err := json.Marshal(d)
// 		if err != nil {
// 			return tempStatusCode, payload, err
// 		}
// 		err = json.Unmarshal(serialised, &deserialised)
// 		if err != nil {
// 			return tempStatusCode, payload, err
// 		}
// 		payload = append(payload, deserialised)
// 	}

// 	return payload, nil
// }
