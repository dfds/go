package model

type Envelope struct {
	Type           string `json:"type"`
	MessageId      string `json:"messageId"`
	EventName      string `json:"eventName"`
	Version        string `json:"version"`
	XCorrelationId string `json:"x-correlationId"`
	XSender        string `json:"x-sender"`
}

type EnvelopeWithPayload[T any] struct {
	Type      string `json:"type"`
	MessageId string `json:"messageId"`
	Payload   T      `json:"data"`
}

type HandlerContext struct {
	Event *Envelope
	Msg   []byte
}
