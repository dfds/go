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

type DoHttpRequestParameters struct {
	HttpClient       *http.Client
	Req              *http.Request
	ParameterSession Session
	DefaultSession   Session
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

type ErrorHandler[T any] interface {
	Handle(e error, response *http.Response) (*T, error)
}

// type DoHTTPRequestResponse struct{}
// T = ErrorResponseEntity that implements Error interface
func DoHttpRequest[T any](req DoHttpRequestParameters, errorHandler ErrorHandler[T]) (*http.Response, *T, error) {
	var errorResponseEntity *T
	err := SetAuthHeader(req.Req, req.ParameterSession, req.DefaultSession)
	if err != nil {
		return nil, errorResponseEntity, err
	}

	resp, err := req.HttpClient.Do(req.Req)
	if err != nil {
		return resp, errorResponseEntity, err
	}

	errorResponseEntity, err = errorHandler.Handle(err, resp)
	if err != nil {
		return nil, errorResponseEntity, err
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

// P = ResponseEntity
// E = ErrorResponseEntity
func HandlePagedResponse[P any, E any](r DoHttpRequestParameters, initialResp Response, errorHandler ErrorHandler[E]) ([]P, *E, error) {
	var errorResponseEntity *E
	payload := []P{}
	container := []interface{}{}
	container = append(container, initialResp.Data)

	nextUrl := initialResp.Metadata.Next

	for {
		if nextUrl == "" {
			break
		}
		req, err := http.NewRequest("GET", nextUrl, nil)
		if err != nil {
			return payload, errorResponseEntity, err
		}
		nextRespRaw, errorResponseEntity, err := DoHttpRequest(DoHttpRequestParameters{
			HttpClient:       r.HttpClient,
			Req:              req,
			ParameterSession: r.ParameterSession,
			DefaultSession:   r.DefaultSession,
		}, errorHandler)
		if err != nil {
			return payload, errorResponseEntity, err
		}
		if errorResponseEntity != nil {
			return payload, errorResponseEntity, err
		}
		nextResp, err := DeserializeResponse[Response](nextRespRaw.Body)
		if err != nil {
			return payload, errorResponseEntity, err
		}
		container = append(container, nextResp.Data)
		nextUrl = nextResp.Metadata.Next
	}
	// Deserialise
	// TODO: Perhaps replace with https://github.com/mitchellh/mapstructure
	denested := []map[string]interface{}{}
	for _, c := range container {
		objects := c.([]interface{})
		for _, object := range objects {
			denested = append(denested, object.(map[string]interface{}))
		}
	}
	for _, d := range denested {
		var deserialised P
		serialised, err := json.Marshal(d)
		if err != nil {
			return payload, errorResponseEntity, err
		}
		err = json.Unmarshal(serialised, &deserialised)
		if err != nil {
			return payload, errorResponseEntity, err
		}
		payload = append(payload, deserialised)
	}
	return payload, errorResponseEntity, nil
}
