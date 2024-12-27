package model

type None struct{}

type ProducerTopic string
type ConsumerTopic string

type ProducerMessage struct {
	Topic     ProducerTopic
	Key       []byte
	Value     []byte
	Partition int32
}

type ProducerError struct {
	Msg *ProducerMessage
	Err error
}
