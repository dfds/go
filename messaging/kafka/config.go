package kafka

type AuthConfig struct {
	Brokers          []string          `required:"true"`
	Mechanism        string            `required:"false"`
	MechanismOptions map[string]string `envconfig:"MECHANISM_OPTIONS"`
	Tls              bool              `required:"true"`
}
