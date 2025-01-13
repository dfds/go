package messaging

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	kafka2 "github.com/segmentio/kafka-go"
	"go.dfds.cloud/messaging/kafka"
	"go.uber.org/zap"
	"sync"
)

type Config struct {
	EnvVarPrefix    string // example: SSU_KAFKA_AUTH
	Wg              *sync.WaitGroup
	Logger          *zap.Logger
	kafkaAuthConfig kafka.AuthConfig
}

type Messaging struct {
	Config  *Config
	Context context.Context
	dialer  *kafka2.Dialer
}

func CreateMessaging() *Messaging {
	return &Messaging{}
}

func (m *Messaging) Init(ctx context.Context, cfg *Config) error {
	m.Config = cfg
	m.Context = ctx

	var authConfig kafka.AuthConfig
	err := envconfig.Process(m.Config.EnvVarPrefix, &authConfig)
	if err != nil {
		return err
	}
	m.Config.kafkaAuthConfig = authConfig

	dialer, err := kafka.NewDialer(m.Config.EnvVarPrefix, authConfig)
	if err != nil {
		return err
	}
	m.dialer = dialer

	return nil
}

func (m *Messaging) NewConsumer(topicName string, groupId string) *kafka.Consumer {
	consumer := kafka.NewConsumer(topicName, groupId, m.Config.kafkaAuthConfig, m.dialer, m.Config.Logger, m.Config.Wg, m.Context)
	return consumer
}

//func StartEventHandling(ctx context.Context, config *Config) error {
//	var authConfig kafka.AuthConfig
//	err := envconfig.Process(config.EnvVarPrefix, &authConfig)
//	if err != nil {
//		return err
//	}
//
//	dialer, err := kafka.NewDialer(authConfig)
//	if err != nil {
//		return err
//	}
//
//	auditConsumer := kafka.NewConsumer("cloudengineering.selfservice.audit", "cloudengineering.ssu-k8s", authConfig, dialer, config.Logger, config.Wg, ctx)
//	auditConsumer.Register("user-action", func(ctx context.Context, event model.HandlerContext) error {
//		fmt.Println("pog")
//		return nil
//	})
//	go auditConsumer.StartConsumer()
//
//	return nil
//}
