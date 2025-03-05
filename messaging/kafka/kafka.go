package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"go.dfds.cloud/messaging/kafka/model"
	"go.dfds.cloud/messaging/kafka/registry"
	"go.uber.org/zap"
	"io"
	"log"
	"sync"
)

func newConsumer(topic string, groupId string, authConfig AuthConfig, dialer *kafka.Dialer) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: authConfig.Brokers,
		GroupID: groupId,
		Topic:   topic,
		Dialer:  dialer,
	})
}

func NewConsumer(topic string, groupId string, authConfig AuthConfig, dialer *kafka.Dialer, logger *zap.Logger, wg *sync.WaitGroup, ctx context.Context) *Consumer {
	return &Consumer{
		topic:      topic,
		groupId:    groupId,
		Reader:     newConsumer(topic, groupId, authConfig, dialer),
		registry:   registry.NewRegistry(),
		logger:     logger,
		wg:         wg,
		ctx:        ctx,
		authConfig: authConfig,
		dialer:     dialer,
	}
}

type Consumer struct {
	topic      string
	groupId    string
	authConfig AuthConfig
	dialer     *kafka.Dialer
	ctx        context.Context
	registry   *registry.Registry
	Reader     *kafka.Reader
	logger     *zap.Logger
	wg         *sync.WaitGroup
}

func (c *Consumer) Topic() string {
	return c.topic
}

func (c *Consumer) Register(eventName string, f registry.HandlerFunc) {
	c.registry.Register(eventName, f)
}

func (c *Consumer) StartConsumer() {
	var cleanupOnce sync.Once
	partitionOffsetTracker := make(map[int]int64)
	cleanup := func() {
		c.logger.Debug("Closing Kafka consumer")
		if err := c.Reader.Close(); err != nil {
			c.logger.Fatal("Failed to close Kafka consumer", zap.Error(err))
		}
		c.logger.Debug("Kafka consumer has been closed")
	}
	defer cleanupOnce.Do(cleanup)

	c.wg.Add(1)
	defer c.wg.Done()

	for {
		c.logger.Debug("Awaiting new message from topic")
		msg, err := c.Reader.FetchMessage(c.ctx)
		if err == io.EOF {
			c.logger.Info("Connection closed")
			break
		} else if err == context.Canceled {
			c.logger.Info("Processing canceled")
			break
		} else if err != nil {
			c.logger.Error("Error fetching message", zap.Error(err))
			break
		}
		msgLog := c.logger.With(zap.String("topic", msg.Topic),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.String("key", string(msg.Key)),
			zap.String("value", string(msg.Value)))
		msgLog.Debug("Message fetched")

		// Convert msg to Event
		var handler registry.HandlerFunc
		var eventLog *zap.Logger

		event, err := GetEventFromMsg(msg.Value)
		if err != nil {
			msgLog.Info("Unable to deserialise message payload. Quite likely the message is not valid JSON. Skipping message")
			partitionOffsetTracker[msg.Partition] = msg.Offset + 1
			continue
		}

		if event == nil || (event.Type == "" && event.EventName == "") {
			msgLog.Info("Unable to recognise event envelope, skipping message")
			partitionOffsetTracker[msg.Partition] = msg.Offset + 1
			continue
		}

		eventLog = c.logger.With(zap.String("topic", msg.Topic),
			zap.Int("partition", msg.Partition),
			zap.Int64("offset", msg.Offset),
			zap.String("key", string(msg.Key)),
			zap.String("eventName", event.Type))

		handlerType := ""
		if event.Type != "" {
			handlerType = event.Type
		} else {
			handlerType = event.EventName
		}
		
		handler = c.registry.GetHandler(handlerType)
		if handler == nil {
			eventLog.Info("No handler registered for event, skipping.")
			partitionOffsetTracker[msg.Partition] = msg.Offset + 1
			continue
		}

		err = handler(c.ctx, model.HandlerContext{
			Event: event,
			Msg:   msg.Value,
		})
		if err != nil {
			eventLog.Error("Handler for event failed", zap.Error(err))
			cleanupOnce.Do(cleanup)
			log.Fatal(err)
		}

		partitionOffsetTracker[msg.Partition] = msg.Offset + 1
	}

	offsets := make(map[string]map[int]int64)
	offsets[c.topic] = partitionOffsetTracker

	c.logger.Info("Updating offsets")
	cg, err := kafka.NewConsumerGroup(kafka.ConsumerGroupConfig{
		ID:      c.groupId,
		Brokers: c.authConfig.Brokers,
		Dialer:  c.dialer,
		Topics:  []string{c.topic},
	})
	if err != nil {
		c.logger.Error("Unable to update commit for consumer group", zap.Error(err)) // TODO: Trigger graceful shutdown
	}

	gen, err := cg.Next(context.Background())
	if err != nil {
		c.logger.Error("Unable to update commit for consumer group", zap.Error(err)) // TODO: Trigger graceful shutdown
	}
	err = gen.CommitOffsets(offsets)
	if err != nil {
		c.logger.Error("Unable to update commit for consumer group", zap.Error(err)) // TODO: Trigger graceful shutdown
	}

	cg.Close()

	cleanupOnce.Do(cleanup)
}

func GetEventFromMsg(data []byte) (*model.Envelope, error) {
	var payload *model.Envelope
	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func commitMsg(ctx context.Context, msg kafka.Message, consumer *kafka.Reader, logger *zap.Logger) error {
	err := consumer.CommitMessages(ctx, msg)
	if err != nil {
		return err
	}

	msgLog := logger.With(zap.String("topic", msg.Topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset))
	msgLog.Debug("Commit for consumer group updated")

	return nil
}
