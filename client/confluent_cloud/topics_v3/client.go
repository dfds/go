package topics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	confluent_util "go.dfds.cloud/client/confluent_cloud/util"
)

const (
	errInvalidRequestFormat = "request does not contain valid data"
)

type TopicsClient struct {
	defaultSession confluent_util.Session
	http           *http.Client
}

func NewClient(defaultSession confluent_util.Session, defaultHttp *http.Client) *TopicsClient {
	return &TopicsClient{
		defaultSession: defaultSession,
		http:           defaultHttp,
	}
}

// Create Request
type CreateTopicsRequest struct {
	Endpoint  string
	ClusterID string
	Payload   CreateTopicsRequestPayload
}

type CreateTopicsRequestPayload struct {
	TopicName         string `json:"topic_name"`
	PartitionsCount   int    `json:"partitions_count,omitempty"`
	ReplicationFactor int    `json:"replication_factor,omitempty"`
	Configs           []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"configs,omitempty"`
}

type CreateTopicResponseEntity struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Self         string `json:"self"`
		ResourceName string `json:"resource_name"`
	} `json:"metadata"`
	ClusterID         string `json:"cluster_id"`
	TopicName         string `json:"topic_name"`
	IsInternal        bool   `json:"is_internal"`
	ReplicationFactor int    `json:"replication_factor"`
	PartitionsCount   int    `json:"partitions_count"`
	Partitions        struct {
		Related string `json:"related"`
	} `json:"partitions"`
	Configs struct {
		Related string `json:"related"`
	} `json:"configs"`
	PartitionReassignments struct {
		Related string `json:"related"`
	} `json:"partition_reassignments"`
	AuthorizedOperations []interface{} `json:"authorized_operations"`
}

func (c *TopicsClient) CreateKafkaTopic(session confluent_util.Session, request CreateTopicsRequest) (CreateTopicResponseEntity, error) {
	var payload CreateTopicResponseEntity
	requestPayload, err := json.Marshal(request.Payload)
	if err != nil {
		return payload, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics"), bytes.NewBuffer(requestPayload))
	if err != nil {
		return payload, err
	}
	payload, err = confluent_util.DoHttpRequest[CreateTopicResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	})
	if err != nil {
		return payload, err
	}
	return payload, nil
}
