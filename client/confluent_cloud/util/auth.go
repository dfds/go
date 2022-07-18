package util

var (
	CLUSTER_API_KEY_SESSION_TYPE = "CLUSTER_API_KEY_SESSION_TYPE"
	CLOUD_API_KEY_SESSION_TYPE   = "CLOUD_API_KEY_SESSION_TYPE"
)

type Session interface {
	GetKey() string
	GetType() string
}

// CloudApiKey

type CloudApiKeySession struct {
	key string
}

func NewCloudApiKeySession(key string) *CloudApiKeySession {
	return &CloudApiKeySession{key: key}
}

func (c *CloudApiKeySession) GetKey() string {
	return c.key
}

func (c *CloudApiKeySession) GetType() string {
	return CLOUD_API_KEY_SESSION_TYPE
}

// ClusterApiKey

type ClusterApiKeySession struct {
	key string
}

func NewClusterApiKeySession(key string) *ClusterApiKeySession {
	return &ClusterApiKeySession{key: key}
}

func (c *ClusterApiKeySession) GetKey() string {
	return c.key
}

func (c *ClusterApiKeySession) GetType() string {
	return CLUSTER_API_KEY_SESSION_TYPE
}
