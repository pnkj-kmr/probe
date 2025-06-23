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

	st := time.Now()
	exitCh := make(chan M.None)

	consumer, err := kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("consumer created")
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println("commence consuming")
	// partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		panic(err)
	}

	go consumer_callback(partitionConsumer, exitCh)

	fmt.Println("message consumed time -- ", time.Since(st))
	<-exitCh
}

func consumer_callback(c sarama.PartitionConsumer, exitCh chan M.None) {
	st := time.Now()
	defer func() {
		if err := c.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-c.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed--- %d\n", consumed)
	fmt.Println("completed message flushed time -- ", time.Since(st))
	exitCh <- M.None{}
}
