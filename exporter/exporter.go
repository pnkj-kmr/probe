package exporter

import (
	"errors"
	"fmt"
	"probe/exporter/kafka"

	"github.com/IBM/sarama"
)

var (
	ErrNoExport = errors.New("NO EXPORT CONFIGURATION")
)

// Exporter
type Exporter struct {
	id    int
	name  string
	kafka sarama.AsyncProducer
	input chan []byte
}

func New(id int, name string) (*Exporter, error) {
	switch name {
	case "kafka":
		kf, err := kafka.NewProducer([]string{"pankaj.local:9092"})
		if err != nil {
			return nil, err
		}
		e := &Exporter{
			id:    id,
			name:  name,
			input: make(chan []byte),
			kafka: kf,
		}
		go e.loop()
		go e.loop2()
		return e, nil
	default:
		return nil, ErrNoExport
	}
}

func (r *Exporter) loop() {
	switch r.name {
	case "kafka":
		for d := range r.input {
			r.kafka.Input() <- &sarama.ProducerMessage{
				Topic: "default_topic", Key: nil, Value: sarama.ByteEncoder(d)}

		}
	default:
	}
}

func (r *Exporter) Export() (data chan []byte) {
	// in case of kafka
	// - need to group the extra metadata into message
	// - multi message should be group together to push to kafka
	// - it will help to partition based topic push
	// - mode information needed like org, agent_id or etc
	fmt.Println("kafka push------", r.name)
	return r.input
}

func (r *Exporter) loop2() {
	// var enqueued int
	var producerErrors, producerSuccess int
	// ProducerLoop:
	for {
		select {
		case _, ok := <-r.kafka.Errors():
			if ok {
				// log.Println("error --- ", err)
				producerErrors++
			}
		case _, ok := <-r.kafka.Successes():
			if ok {
				// log.Println("pushed------", string(msg.Value))
				producerSuccess++
			}
		}
	}
	// fmt.Println("success--", producerSuccess, "error----", producerErrors)
}
