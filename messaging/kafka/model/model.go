package model

import "github.com/segmentio/kafka-go"

type Envelope struct {
	Type           string `json:"type"`
	MessageId      string `json:"messageId"`
	EventName      string `json:"eventName"`
	Version        string `json:"version"`
	XCorrelationId string `json:"x-correlationId"`
	XSender        string `json:"x-sender"`
}

type EnvelopeWithPayload[T any] struct {
	Type           string `json:"type,omitempty"`
	MessageId      string `json:"messageId,omitempty"`
	EventName      string `json:"eventName"`
	Version        string `json:"version"`
	XCorrelationId string `json:"x-correlationId"`
	XSender        string `json:"x-sender"`
	Payload        T      `json:"payload"`
}

type HandlerContext struct {
	Event  *Envelope
	Msg    []byte
	Writer NewWriterFunc
}

type NewWriterFunc func(topic string) *kafka.Writer
