package confluent_cloud

import (
	"net/http"

	"go.dfds.cloud/client/confluent_cloud/service_accounts_v2"
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

func (c *Client) ServiceAccountsV2() *service_accounts_v2.ServiceAccountsClient {
	return service_accounts_v2.NewClient(c.defaultSession, c.http)
}

// func (c *Client) TopicsV3() *topics.TopicsClient {
// 	return topics.NewClient(c.defaultSession, c.http, c.defaultEndpoint)
// }

// func (c *Client) TopicsV3() *topics.TopicsClient {
// 	return topics.NewClient(c.defaultSession, c.http)
// }

type KafkaClient struct {
	defaultSession  confluent_util.Session
	http            *http.Client
	defaultEndpoint *string
	defaultCluster  *string
}

func NewKafkaClient() *KafkaClient {
	return &KafkaClient{
		defaultSession:  nil,
		http:            http.DefaultClient,
		defaultEndpoint: nil,
		defaultCluster:  nil,
	}
}

func (c *KafkaClient) SetKafkaClientDefaultClusterInfo(endpoint string, id string, session confluent_util.Session) {
	c.defaultEndpoint = &endpoint
	c.defaultCluster = &id
	c.defaultSession = session
}

func (c *KafkaClient) TopicsV3() *topics.TopicsClient {
	return topics.NewClient(c.defaultSession, c.http, c.defaultEndpoint, c.defaultCluster)
}
