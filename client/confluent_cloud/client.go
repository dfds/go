package confluent_cloud

import (
	"net/http"

	topics "go.dfds.cloud/client/confluent_cloud/topics_v3"
	confluent_util "go.dfds.cloud/client/confluent_cloud/util"
)

type Client struct {
	defaultSession confluent_util.Session
	http           *http.Client
}

func NewClient() *Client {
	return &Client{
		defaultSession: nil,
		http:           http.DefaultClient,
	}
}

func (c *Client) Authenticate() error {
	return nil
}

func (c *Client) SetDefaultSession(session confluent_util.Session) {
	c.defaultSession = session
}

// func (c *Client) ServiceAccountsV2() *service_accounts_v2.ServiceAccountsClient {
// 	return service_accounts_v2.NewClient(c.defaultSession, c.http)
// }

func (c *Client) TopicsV3() *topics.TopicsClient {
	return topics.NewClient(c.defaultSession, c.http)
}
