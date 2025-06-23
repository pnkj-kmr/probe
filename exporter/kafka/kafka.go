package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

/**
Main configration for kafka related

`json:"version"`
`json:"bootstrap.servers"`
`json:"security.protocol"`
`json:"ssl.ca.location"`
`json:"ssl.endpoint.identification.algorithm"`
`json:"sasl.mechanism"`
`json:"sasl.username"`
`json:"sasl.password"`
`json:"message.max.bytes"`
`json:"queue.buffering.max.messages"`
`json:"retries"`
`json:"linger.ms"`
`json:"batch.num.messages"`
`json:"partitioner"`
`json:"acks,omitempty"`
`json:"compression.type"`

`json:"client.id"`
`json:"auto.offset.reset"`
`json:"group.id"`
`json:"session.timeout.ms"`

*/

func DefaultProducerConfig() (*sarama.Config, error) {
	conf := sarama.NewConfig()

	// TODO - need to handle
	// version, err := sarama.ParseKafkaVersion(pc.Version)
	// sarama.DefaultVersion.String()
	// if err != nil {
	// 	return nil, err
	// }
	conf.Version = sarama.V3_0_1_0
	conf.ClientID = "probe" // TODO
	// conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "infraon-producer"
	conf.Net.SASL.Password = "infraon@123"
	// conf.Net.SASL.Handshake = true
	sasl := "SCRAM-SHA-512"
	if sasl == "SCRAM-SHA-512" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	} else if sasl == "SCRAM-SHA-256" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	} else {
		return nil, ErrNoSASLMechanismFound
	}
	// certFile      = flag.String("certificate", "", "The optional certificate file for client authentication")
	// keyFile       = flag.String("key", "", "The optional key file for client authentication")
	// caFile        = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	caFile := "ca-cert"
	// TODO - tls flag - default true
	t, err := createTLSConfiguration(true, caFile, "", "")
	if err != nil {
		return nil, err
	}
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = t

	//-------------------------------------------------
	// sarama.MaxRequestSize
	// conf.Producer.Compression = sarama.CompressionNone
	conf.Producer.Compression = sarama.CompressionSnappy
	conf.Producer.MaxMessageBytes = 20 * 1024 * 1024
	conf.Producer.Retry.Max = 3

	// conf.Producer.Flush = 10 * 1024 * 1024
	// conf.Producer.Timeout = time.Second * 6
	conf.Producer.Flush.Frequency = time.Second * 1
	conf.Producer.Flush.Messages = 10
	conf.Producer.Flush.Bytes = 20 * 1024 * 1024
	conf.Producer.Flush.MaxMessages = 1000

	conf.Producer.RequiredAcks = sarama.WaitForAll // p.Acks
	conf.Producer.Return.Successes = true          // default false
	conf.Producer.Return.Errors = true             // default true

	// final validation of kafka producer config
	err = conf.Validate()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func DefaultConsumerConfig() (*sarama.Config, error) {
	conf := sarama.NewConfig()

	//	    'sasl.username': 'dnsadmin',
	//	    'sasl.password': 'dnsadmin@123',
	//	    'group.id': 'consumer-active-events',
	//	    'session.timeout.ms': 6000,
	//	    'auto.offset.reset': 'earliest',
	//	    'client.id': "consumer-active-events-%s" % instanceID,
	//	    "message.max.bytes": "10000000"
	// conf.Consumer.Group.InstanceId = "test-group"
	// conf.Consumer.Offsets.Initial = sarama.OffsetOldest

	conf.Version = sarama.V3_0_1_0
	conf.ClientID = "probe" // TODO
	// conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = "dnsadmin"
	conf.Net.SASL.Password = "dnsadmin@123"
	// conf.Net.SASL.Handshake = true
	sasl := "SCRAM-SHA-512"
	if sasl == "SCRAM-SHA-512" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	} else if sasl == "SCRAM-SHA-256" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	} else {
		return nil, ErrNoSASLMechanismFound
	}
	// certFile      = flag.String("certificate", "", "The optional certificate file for client authentication")
	// keyFile       = flag.String("key", "", "The optional key file for client authentication")
	// caFile        = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	caFile := "ca-cert"
	// TODO - tls flag - default true
	t, err := createTLSConfiguration(true, caFile, "", "")
	if err != nil {
		return nil, err
	}
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = t

	// final validate
	err = conf.Validate()
	if err != nil {
		return nil, err
	}
	return conf, nil
}
