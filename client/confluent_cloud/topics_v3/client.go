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

type ErrorResponseEntity struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func (e *ErrorResponseEntity) Error() string {
	return e.Message
}

// T = Payload
// E = any or interface{}
type TopicRequestEntity[T any, E any] struct { // Accepts both Payload and Error
	Endpoint  string
	ClusterID string
	TopicName *string `json:"topic_name,omitempty"`
	Payload   *T
}

func (c *TopicRequestEntity[T, E]) Handle(e error, r *http.Response) (*E, error) { // Will be called automatically on CreateTopicsRequest - implements the interface
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

type CreateTopicsRequestPayload struct {
	TopicName         string              `json:"topic_name"`
	PartitionsCount   int                 `json:"partitions_count,omitempty"`
	ReplicationFactor int                 `json:"replication_factor,omitempty"`
	Configs           []TopicConfigEntity `json:"configs,omitempty"`
}

type CreateTopicResponsePayload struct {
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

func (c *TopicsClient) CreateKafkaTopic(session confluent_util.Session, request TopicRequestEntity[CreateTopicsRequestPayload, ErrorResponseEntity]) (CreateTopicResponsePayload, error) {
	var payload CreateTopicResponsePayload
	requestPayload, err := json.Marshal(request.Payload)
	if err != nil {
		return payload, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics"), bytes.NewBuffer(requestPayload))
	if err != nil {
		return payload, err
	}
	r, errorResponseEntity, err := confluent_util.DoHttpRequest[ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	},
		&request,
	)

	if err != nil {
		return payload, err
	}

	if errorResponseEntity != nil {
		return payload, errorResponseEntity
	}

	payload, err = confluent_util.DeserializeResponse[CreateTopicResponsePayload](r.Body)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

// DELETE request
func (c *TopicsClient) DeleteKafkaTopic(session confluent_util.Session, request TopicRequestEntity[any, ErrorResponseEntity]) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics", *request.TopicName), nil)
	if err != nil {
		return err
	}
	_, errorResponseEntity, err := confluent_util.DoHttpRequest[ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, &request)
	if err != nil {
		return err
	}
	if errorResponseEntity != nil {
		return errorResponseEntity
	}
	return nil
}

// Get request
type GetTopicResponseEntity struct { // Note: same as CreateTopicResponsePayload
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
}

func (c *TopicsClient) GetKafkaTopic(session confluent_util.Session, request TopicRequestEntity[any, ErrorResponseEntity]) (GetTopicResponseEntity, error) {
	var payload GetTopicResponseEntity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics", *request.TopicName), nil)
	if err != nil {
		return payload, err
	}
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
	payload, err = confluent_util.DeserializeResponse[GetTopicResponseEntity](r.Body)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

type ListTopicResponseEntity struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Self string      `json:"self"`
		Next interface{} `json:"next"`
	} `json:"metadata"`
	Data []GetTopicResponseEntity `json:"data"`
}

func (c *TopicsClient) ListKafkaTopic(session confluent_util.Session, request TopicRequestEntity[any, ErrorResponseEntity]) (ListTopicResponseEntity, error) {
	var payload ListTopicResponseEntity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics"), nil)
	if err != nil {
		return payload, err
	}
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

	payload, err = confluent_util.DeserializeResponse[ListTopicResponseEntity](r.Body)
	if err != nil {
		return payload, err
	}

	return payload, nil
}

type TopicConfigEntity struct {
	Name      string `json:"name"`
	Value     string `json:"value,omitempty"`
	Operation string `json:"operation,omitempty"`
}

type UpdateTopicConfigRequestPayload struct {
	Data []TopicConfigEntity `json:"data"`
}

type UpdateTopicsRequest struct {
	Endpoint  string
	ClusterID string
	TopicName string
	Payload   UpdateTopicConfigRequestPayload
}

func (c *TopicsClient) UpdateKafkaTopic(session confluent_util.Session, request TopicRequestEntity[UpdateTopicConfigRequestPayload, ErrorResponseEntity]) error {
	requestPayload, err := json.Marshal(request.Payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics", *request.TopicName, "configs:alter"), bytes.NewBuffer(requestPayload))
	if err != nil {
		return err
	}

	_, errorResponseEntity, err := confluent_util.DoHttpRequest[ErrorResponseEntity](confluent_util.DoHttpRequestParameters{
		HttpClient:       c.http,
		Req:              req,
		ParameterSession: session,
		DefaultSession:   c.defaultSession,
	}, &request)
	if err != nil {
		return err
	}
	if errorResponseEntity != nil {
		return errorResponseEntity
	}
	return nil
}

type GetTopicConfigsResponseEntity struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Self string      `json:"self"`
		Next interface{} `json:"next"`
	} `json:"metadata"`
	Data []GetTopicConfigResponseEntity `json:"data"`
}

type GetTopicConfigResponseEntity struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Self         string `json:"self"`
		ResourceName string `json:"resource_name"`
	} `json:"metadata"`
	ClusterID   string `json:"cluster_id"`
	TopicName   string `json:"topic_name"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	IsDefault   bool   `json:"is_default"`
	IsReadOnly  bool   `json:"is_read_only"`
	IsSensitive bool   `json:"is_sensitive"`
	Source      string `json:"source"`
	Synonyms    []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Source string `json:"source"`
	} `json:"synonyms"`
}

func (c *TopicsClient) GetKafkaTopicConfigs(session confluent_util.Session, request TopicRequestEntity[any, ErrorResponseEntity]) (GetTopicConfigsResponseEntity, error) {
	var payload GetTopicConfigsResponseEntity
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s/%s/%s/%s", request.Endpoint, "kafka/v3/clusters", request.ClusterID, "topics", *request.TopicName, "configs"), nil)
	if err != nil {
		return payload, err
	}
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
	payload, err = confluent_util.DeserializeResponse[GetTopicConfigsResponseEntity](r.Body)
	if err != nil {
		return payload, err
	}

	return payload, nil
}
