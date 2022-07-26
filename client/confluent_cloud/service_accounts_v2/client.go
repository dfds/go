package service_accounts_v2

import (
	"fmt"
	"net/http"
	"time"

	"go.dfds.cloud/client/confluent_cloud/util"
	confluent_util "go.dfds.cloud/client/confluent_cloud/util"
)

type ServiceAccountsClient struct {
	defaultSession confluent_util.Session
	http           *http.Client
}

func NewClient(defaultSession confluent_util.Session, defaultHttp *http.Client) *ServiceAccountsClient {
	return &ServiceAccountsClient{
		defaultSession: defaultSession,
		http:           defaultHttp,
	}
}

type ErrorResponseEntity struct {
	Errors []struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}

func (e *ErrorResponseEntity) Error() string { // TODO: Fix this because it should read from ErrorResponseEntity and not to always return string!!
	// fmt.Println(e.Errors)
	return "Some Error" // e.Message
}

type ServiceAccountRequestEntity[T any, E any] struct {
	ServiceAccountName *string `json:"service_account_name,omitempty"`
	Payload            *T
}

func (c *ServiceAccountRequestEntity[T, E]) Handle(e error, r *http.Response) (*E, error) { // Will be called automatically on CreateTopicsRequest - implements the interface
	var errorMessage *E
	if r.StatusCode != 204 && r.StatusCode >= 400 {
		errorMessage, err := confluent_util.DeserializeResponse[E](r.Body)
		if err != nil {
			return nil, err
		}
		return &errorMessage, nil
	}
	return errorMessage, nil
}

type PaginationParameters struct {
	PageSize  string `json:"page_size"`
	PageToken string `json:"page_token"`
}

// type GetServiceAccountsRequest struct {
// 	PageSize  string `json:"page_size"`
// 	PageToken string `json:"page_token"`
// }

type GetServiceAccountResponseEntity struct {
	APIVersion  string `json:"api_version"`
	Description string `json:"description"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Metadata    struct {
		CreatedAt    time.Time `json:"created_at"`
		ResourceName string    `json:"resource_name"`
		Self         string    `json:"self"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"metadata"`
}

func (c *ServiceAccountsClient) ListServiceAccounts(session confluent_util.Session, request ServiceAccountRequestEntity[any, ErrorResponseEntity],
	requestQueryParameters PaginationParameters) ([]GetServiceAccountResponseEntity, error) {
	var payload []GetServiceAccountResponseEntity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", confluent_util.BASE_ENDPOINT, "/iam/v2/service-accounts"), nil)
	if err != nil {
		return payload, err
	}

	query := req.URL.Query()
	query.Add("page_size", requestQueryParameters.PageSize)
	query.Add("page_token", requestQueryParameters.PageToken)
	req.URL.RawQuery = query.Encode()

	r, errorResponseEntity, err := confluent_util.DoHttpRequest[ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, &request)
	if err != nil {
		fmt.Println("Error line 100")
		return payload, err
	}
	if errorResponseEntity != nil {
		fmt.Println("Error line 104")
		return payload, errorResponseEntity
	}
	initialResp, err := confluent_util.DeserializeResponse[util.Response](r.Body)
	if err != nil {
		fmt.Println("Error line 109")
		return payload, err
	}
	// ==

	payload, errorResponseEntity, err = confluent_util.HandlePagedResponse[GetServiceAccountResponseEntity, ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              nil,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, initialResp, &request)
	if err != nil {
		fmt.Println("Error line 120")
		return payload, err
	}
	if errorResponseEntity != nil {
		fmt.Println("Error line 124")
		return payload, errorResponseEntity
	}
	// ===

	// bodyBytes, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	return payload, err
	// }
	// fmt.Println(string(bodyBytes))
	fmt.Println("Client is done?")
	return payload, err // TODO: This should return errorResponseEntity!! and Test error message with error from client
}

// func (c *ServiceAccountsClient) GetServiceAccounts(session confluent_util.Session, request GetServiceAccountsRequest) ([]GetServiceAccountResponseEntity, error) {
// 	var payload []GetServiceAccountResponseEntity
// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", confluent_util.BASE_ENDPOINT, "/iam/v2/service-accounts"), nil)
// 	if err != nil {
// 		return payload, err
// 	}

// 	query := req.URL.Query()
// 	query.Add("page_size", request.PageSize)
// 	query.Add("page_token", request.PageToken)
// 	req.URL.RawQuery = query.Encode()

// 	resp, err := confluent_util.DoHttpRequest[confluent_util.Response](confluent_util.DoHttpRequestParameters{
// 		HttpClient:       c.http,
// 		Req:              req,
// 		ParameterSession: session,
// 		DefaultSession:   c.defaultSession,
// 	})
// 	if err != nil {
// 		return payload, err
// 	}

// 	payload, err = confluent_util.HandlePagedResponse[GetServiceAccountResponseEntity](confluent_util.DoHttpRequestParameters{
// 		HttpClient:       c.http,
// 		Req:              nil,
// 		ParameterSession: session,
// 		DefaultSession:   c.defaultSession,
// 	}, resp)
// 	if err != nil {
// 		return payload, err
// 	}
// 	return payload, err
// }
