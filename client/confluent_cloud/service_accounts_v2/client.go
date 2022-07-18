package service_accounts_v2

import (
	"fmt"
	confluent_util "go.dfds.cloud/client/confluent_cloud/util"
	"net/http"
	"time"
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

type GetServiceAccountsRequest struct {
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

func (c *ServiceAccountsClient) GetServiceAccounts(session confluent_util.Session, request GetServiceAccountsRequest) ([]GetServiceAccountResponseEntity, error) {
	var payload []GetServiceAccountResponseEntity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", confluent_util.BASE_ENDPOINT, "/iam/v2/service-accounts"), nil)
	if err != nil {
		return payload, err
	}

	query := req.URL.Query()
	query.Add("page_size", request.PageSize)
	query.Add("page_token", request.PageToken)
	req.URL.RawQuery = query.Encode()

	resp, err := confluent_util.DoHttpRequest(confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	})
	if err != nil {
		return payload, err
	}

	payload, err = confluent_util.HandlePagedResponse[GetServiceAccountResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              nil,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, resp)
	if err != nil {
		return payload, err
	}
	return payload, err
}
