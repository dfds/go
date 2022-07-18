package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func DoHttpRequest(req DoHttpRequestParameters) (Response, error) {
	var baseResp Response
	err := SetAuthHeader(req.Req, req.ParameterSession, req.DefaultSession)
	if err != nil {
		return baseResp, err
	}

	resp, err := req.HttpClient.Do(req.Req)
	if err != nil {
		return baseResp, err
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return baseResp, err
	}

	err = json.Unmarshal(respBytes, &baseResp)
	if err != nil {
		return baseResp, err
	}

	return baseResp, nil
}

func HandlePagedResponse[T any](r DoHttpRequestParameters, initialResp Response) ([]T, error) {
	payload := []T{}
	container := []interface{}{}
	container = append(container, initialResp.Data)

	nextUrl := initialResp.Metadata.Next

	for {
		if nextUrl == "" {
			break
		}
		fmt.Println(nextUrl)
		req, err := http.NewRequest("GET", nextUrl, nil)
		if err != nil {
			return payload, err
		}

		nextResp, err := DoHttpRequest(DoHttpRequestParameters{
			HttpClient:       r.HttpClient,
			Req:              req,
			ParameterSession: r.ParameterSession,
			DefaultSession:   r.DefaultSession,
		})
		if err != nil {
			return payload, err
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
		var deserialised T
		serialised, err := json.Marshal(d)
		if err != nil {
			return payload, err
		}
		err = json.Unmarshal(serialised, &deserialised)
		if err != nil {
			return payload, err
		}
		payload = append(payload, deserialised)
	}

	return payload, nil
}
