package kafka

import (
	"github.com/IBM/sarama"
)

// type Producer interface {
// 	sarama.AsyncProducer
// 	Produce([]byte)
// }
// type producer struct {
// 	isClosed bool
// 	done     chan M.None
// 	producer sarama.AsyncProducer
// 	// errors           chan *M.ProducerError
// 	// input, successes chan *M.ProducerMessage
// 	// // TODO - to bind extra metadata or log
// }

func NewProducer(bootstrapServers []string) (sarama.AsyncProducer, error) {
	conf, err := DefaultProducerConfig()
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewAsyncProducer(bootstrapServers, conf)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// func (p *producer) Produce([]byte) {

// }

// func (p *producer) AsyncClose()                               { p.producer.AsyncClose() }
// func (p *producer) Close() error                              { return p.producer.Close() }
// func (p *producer) Input() chan<- *sarama.ProducerMessage     { return p.Input() }
// func (p *producer) Successes() <-chan *sarama.ProducerMessage { return p.Successes() }
// func (p *producer) Errors() <-chan *sarama.ProducerError      { return p.Errors() }
// func (p *producer) IsTransactional() bool                     { return p.producer.IsTransactional() }
// func (p *producer) TxnStatus() sarama.ProducerTxnStatusFlag   { return p.producer.TxnStatus() }
// func (p *producer) BeginTxn() error                           { return p.producer.BeginTxn() }
// func (p *producer) CommitTxn() error                          { return p.producer.CommitTxn() }
// func (p *producer) AbortTxn() error                           { return p.producer.AbortTxn() }
// func (p *producer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
// 	return p.producer.AddOffsetsToTxn(offsets, groupId)
// }
// func (p *producer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
// 	return p.producer.AddMessageToTxn(msg, groupId, metadata)
// }
