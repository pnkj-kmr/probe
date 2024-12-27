package kafka

import "probe/model"

type Producer interface {
	IsClosed() bool
	Done() chan<- model.None
	Produce() chan<- *model.ProducerMessage
	OnSuccess() <-chan *model.ProducerMessage
	OnError() <-chan *model.ProducerError
}

type Config struct {
	Version              string   `json:"version"`
	BootstrapServers     []string `json:"bootstrap.servers"`
	SecurityProtocol     string   `json:"security.protocol"`
	SSLCAlocation        string   `json:"ssl.ca.location"`
	SSLEndpointAlgorithm string   `json:"ssl.endpoint.identification.algorithm"`
	SASLMechanism        string   `json:"sasl.mechanism"`
	SASLUsername         string   `json:"sasl.username"`
	SASLPassword         string   `json:"sasl.password"`
	MessageMaxBytes      string   `json:"message.max.bytes"`
}

type ProducerConfig struct {
	*Config
	// producer related config
	QueueBufferMaxMessages int    `json:"queue.buffering.max.messages"`
	Retires                int    `json:"retries"`
	LingerMS               int    `json:"linger.ms"`
	BatchNumMessages       int    `json:"batch.num.messages"`
	Partitioner            string `json:"partitioner"`
	Acks                   bool   `json:"acks,omitempty"`
	CompressionType        bool   `json:"compression.type"`
}

type ConsumerConfig struct {
	*Config
	// consumer related config
	ClientID         string `json:"client.id"`
	AutoOffsetReset  string `json:"auto.offset.reset"`
	GroupID          string `json:"group.id"`
	SessionTimeoutMS string `json:"session.timeout.ms"`
}
