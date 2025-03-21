package kafka

import "github.com/IBM/sarama"

/**
TODO - consumer group need to be implemented
	// sarama.ConsumerGroup


*/

func NewConsumer(bootstrapServers []string) (sarama.Consumer, error) {
	conf, err := DefaultConsumerConfig()
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(bootstrapServers, conf)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
