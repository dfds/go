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

func (e *ErrorResponseEntity) Error() string {
	return fmt.Sprint(e.Errors)
}

// T = Payload
// E = any or interface{}
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
		return payload, err
	}
	if errorResponseEntity != nil {
		return payload, errorResponseEntity
	}
	initialResp, err := confluent_util.DeserializeResponse[util.Response](r.Body)
	if err != nil {
		return payload, err
	}

	payload, errorResponseEntity, err = confluent_util.HandlePagedResponse[GetServiceAccountResponseEntity, ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              nil,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, initialResp, &request)
	if err != nil {
		return payload, err
	}
	if errorResponseEntity != nil {
		return payload, errorResponseEntity
	}
	return payload, nil
}
