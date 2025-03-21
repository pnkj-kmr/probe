package main

import (
	"crypto/sha256"
	"crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"probe/exporter/kafka"
	M "probe/model"

	"github.com/IBM/sarama"
	"github.com/xdg-go/scram"
)

var (
	SHA256 scram.HashGeneratorFcn = sha256.New
	SHA512 scram.HashGeneratorFcn = sha512.New
)

type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

func init() {
	sarama.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
}

var (
	brokers       = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	version       = flag.String("version", sarama.DefaultVersion.String(), "Kafka cluster version")
	userName      = flag.String("username", "", "The SASL username")
	passwd        = flag.String("passwd", "", "The SASL password")
	algorithm     = flag.String("algorithm", "", "The SASL SCRAM SHA algorithm sha256 or sha512 as mechanism")
	topic         = flag.String("topic", "default_topic", "The Kafka topic to use")
	certFile      = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile       = flag.String("key", "", "The optional key file for client authentication")
	caFile        = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	tlsSkipVerify = flag.Bool("tls-skip-verify", false, "Whether to skip TLS server cert verification")
	useTLS        = flag.Bool("tls", false, "Use TLS to communicate with the cluster")
	mode          = flag.String("mode", "produce", "Mode to run in: \"produce\" to produce, \"consume\" to consume")
	logMsg        = flag.Bool("logmsg", false, "True to log consumed messages to console")
	msgCount      = flag.Int("c", 10, "default msg write count")

	logger = log.New(os.Stdout, "[Producer] ", log.LstdFlags)
)

func createTLSConfiguration() (t *tls.Config) {
	t = &tls.Config{
		InsecureSkipVerify: *tlsSkipVerify,
	}
	if *certFile != "" && *keyFile != "" && *caFile != "" {
		cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}

		caCert, err := os.ReadFile(*caFile)
		if err != nil {
			log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: *tlsSkipVerify,
		}
	}
	return t
}

func main() {
	flag.Parse()

	if *brokers == "" {
		log.Fatalln("at least one broker is required")
	}
	splitBrokers := strings.Split(*brokers, ",")

	v, err := sarama.ParseKafkaVersion(*version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}

	if *userName == "" {
		log.Fatalln("SASL username is required")
	}

	if *passwd == "" {
		log.Fatalln("SASL password is required")
	}

	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = 1
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.Version = v
	conf.ClientID = "sasl_scram_client"
	conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = *userName
	conf.Net.SASL.Password = *passwd
	conf.Net.SASL.Handshake = true
	if *algorithm == "sha512" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	} else if *algorithm == "sha256" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256

	} else {
		log.Fatalf("invalid SHA algorithm \"%s\": can be either \"sha256\" or \"sha512\"", *algorithm)
	}

	if *useTLS {
		conf.Net.TLS.Enable = true
		conf.Net.TLS.Config = createTLSConfiguration()
	}

	if *mode == "consume" {
		consumer, err := sarama.NewConsumer(splitBrokers, conf)
		if err != nil {
			panic(err)
		}
		log.Println("consumer created")
		defer func() {
			if err := consumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		log.Println("commence consuming")
		partitionConsumer, err := consumer.ConsumePartition(*topic, 0, sarama.OffsetOldest)
		if err != nil {
			panic(err)
		}

		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.Fatalln(err)
			}
		}()

		// Trap SIGINT to trigger a shutdown.
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		consumed := 0
	ConsumerLoop:
		for {
			log.Println("in the for")
			select {
			case msg := <-partitionConsumer.Messages():
				log.Printf("Consumed message offset %d\n", msg.Offset)
				if *logMsg {
					log.Printf("KEY: %s VALUE: %s", msg.Key, msg.Value)
				}
				consumed++
			case <-signals:
				break ConsumerLoop
			}
		}

		log.Printf("Consumed: %d\n", consumed)

	} else {
		// syncProducer, err := sarama.NewSyncProducer(splitBrokers, conf)
		// if err != nil {
		// 	logger.Fatalln("failed to create producer: ", err)
		// }
		// partition, offset, err := syncProducer.SendMessage(&sarama.ProducerMessage{
		// 	Topic: *topic,
		// 	Value: sarama.StringEncoder("hello to kafka from PANKAJ"),
		// })
		// if err != nil {
		// 	logger.Fatalln("failed to send message to ", *topic, err)
		// }
		// logger.Printf("wrote message at partition: %d, offset: %d", partition, offset)
		// _ = syncProducer.Close()

		st := time.Now()
		conf := &kafka.ProducerConfig{
			Config: &kafka.Config{
				Version:          *version,
				BootstrapServers: splitBrokers,
				SSLCAlocation:    *caFile,
				SASLUsername:     *userName,
				SASLPassword:     *passwd,
				SASLMechanism:    "SCRAM-SHA-512",
			},
			Retires: 3,
		}
		p, err := kafka.NewProducer(conf)
		if err != nil {
			log.Fatal("error ----", err)
		}

		// lunching callback
		exitCh := make(chan M.None)
		go producer_callback(p, exitCh)
		for i := 0; i < *msgCount; i++ {
			p.Produce() <- &M.ProducerMessage{
				Topic: M.ProducerTopic(*topic),
				Key:   nil,
				Value: sarama.ByteEncoder(fmt.Sprintf("testing 123 --- %d", i)),
			}
			log.Println("input ----", i)
		}

		p.Done() <- M.None{}
		fmt.Println("message flushed time -- ", time.Since(st))

		<-exitCh

	}
	logger.Println("Bye now !")
}

func producer_callback(p kafka.Producer, exitCh chan M.None) {
	st := time.Now()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// var enqueued int
	var producerErrors, producerSuccess int
ProducerLoop:
	for {
		// log.Println("------", enqueued)
		// if p.IsClosed() {
		// 	log.Println("producer closed...")
		// 	break ProducerLoop
		// }
		select {
		// case p.Produce() <- &M.ProducerMessage{
		// 	Topic: M.ProducerTopic(*topic),
		// 	Key:   nil,
		// 	Value: sarama.ByteEncoder(fmt.Sprintf("testing 123 --- %d", enqueued)),
		// }:
		// 	enqueued++
		// 	if enqueued == *msgCount {
		// 		p.Done() <- M.None{}
		// 	}
		// 	log.Println("input---------", enqueued)
		case _, ok := <-p.OnError():
			if ok {
				// log.Println("error --- ", err)
				producerErrors++
			}
		case _, ok := <-p.OnSuccess():
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
