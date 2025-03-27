package exporter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"probe/exporter/kafka"
	M "probe/model"

	"github.com/IBM/sarama"
)

var (
	ErrNoExport = errors.New("NO EXPORT CONFIGURATION")
)

// Exporter
type Exporter struct {
	ctx         context.Context
	input       chan any
	id          int
	name        string
	kafka       sarama.AsyncProducer
	kfTopic     string
	kfPartition int32
}

func New(ctx context.Context, id int, name string) (*Exporter, error) {
	//
	/*
		TODO
		KAFKA CONFIGURATION UPDATE
	*/
	switch name {
	case "kafka":
		kf, err := kafka.NewProducer([]string{"pankaj.local:9092"})
		if err != nil {
			return nil, err
		}
		e := &Exporter{
			ctx:         ctx,
			id:          id,
			name:        name,
			input:       make(chan any),
			kafka:       kf,
			kfTopic:     "default_topic",
			kfPartition: 4,
		}
		return e, nil
	default:
		return nil, ErrNoExport
	}
}

func (r *Exporter) Spin() {
	go r.loop()
	go r.loop2()
}

func (r *Exporter) Close() {
	if r.kafka != nil {
		r.kafka.AsyncClose()
	}
}

func (r *Exporter) Send() (data chan any) {
	// in case of kafka
	// - need to group the extra metadata into message
	// - multi message should be group together to push to kafka
	// - it will help to partition based topic push
	// - mode information needed like org, agent_id or etc
	fmt.Println("kafka push------", r.name)
	return r.input
}

func (r *Exporter) loop() {
	// TODO loop and loop2 need to be exist gracefully
	//
	switch r.name {
	case "kafka":
		for d := range r.input {
			log.Println("loop1 running.....")
			// if r.kafka == nil {
			// 	return
			// }
			switch msg := d.(type) {
			case M.ExportMsg:
				r.kafka.Input() <- &sarama.ProducerMessage{
					Topic: r.kfTopic,
					Key:   nil, Value: msg,
					Partition: msg.PartitionId,
				}
			default:
				fmt.Println("UNKNOWN MESSAGE RECEIVED")
			}
		}
	default:
	}
}

func (r *Exporter) loop2() {
	// var enqueued int
	var producerErrors, producerSuccess int
	// ProducerLoop:
	for {
		log.Println("loop2 running.....")
		// if r.kafka == nil {
		// 	break
		// }
		select {
		case <-r.ctx.Done():
			break
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
