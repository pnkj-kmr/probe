package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"probe/exporter/kafka"
	M "probe/model"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	topic := "default_topic"
	brokers := []string{"pankaj.local:9092"}
	maxMessages := 1000000

	st := time.Now()
	exitCh := make(chan M.None)

	p, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalln("err---", err)
	}
	go producer_callback(p, exitCh)

	for i := 0; i < maxMessages; i++ {
		p.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Key:   nil,
			Value: sarama.ByteEncoder(fmt.Sprintf("testing 123 --- %d", i)),
		}
		log.Println("input ----", i)
		log.Println("Bye now !")
	}
	fmt.Println("message flushed time -- ", time.Since(st))
	<-exitCh
}

func producer_callback(p sarama.AsyncProducer, exitCh chan M.None) {
	st := time.Now()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// var enqueued int
	var producerErrors, producerSuccess int
ProducerLoop:
	for {
		select {
		case _, ok := <-p.Errors():
			if ok {
				// log.Println("error --- ", err)
				producerErrors++
			}
		case _, ok := <-p.Successes():
			if ok {
				// log.Println("pushed------", string(msg.Value))
				producerSuccess++
			}
		case <-signals:
			log.Println("os kill------")
			break ProducerLoop
		}
	}
	fmt.Println("success--", producerSuccess, "error----", producerErrors)
	fmt.Println("completed message flushed time -- ", time.Since(st))
	exitCh <- M.None{}
}
