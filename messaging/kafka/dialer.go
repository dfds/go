package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"strings"
	"time"
)

func NewDialer(envPrefix string, authConfig AuthConfig) (*kafka.Dialer, error) {
	// Configure TLS

	var tlsConfig *tls.Config
	var saslMechanism sasl.Mechanism

	if authConfig.Tls {
		tlsConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
	}

	// Configure SASL mechanism
	for k, v := range authConfig.MechanismOptions {
		fmt.Printf("Key: %s, Value: %s\n", k, v)
	}

	fmt.Println(authConfig.Mechanism)
	switch strings.ToLower(authConfig.Mechanism) {
	case "plain":
		var plainConf plain.Mechanism
		err := envconfig.Process(fmt.Sprintf("%s_MECHANISM_PLAIN", envPrefix), &plainConf)
		if err != nil {
			return nil, err
		}
		saslMechanism = plainConf
	default:
		saslMechanism = nil
	}

	fmt.Println(saslMechanism)

	// Configure connection dialer
	return &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig,
		SASLMechanism: saslMechanism,
		ClientID:      "ssu-k8s",
	}, nil
}
