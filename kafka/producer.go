package kafka

import (
	"probe/model"

	"github.com/IBM/sarama"
)

type producer struct {
	isClosed bool
	done     chan model.None

	producer         sarama.AsyncProducer
	errors           chan *model.ProducerError
	input, successes chan *model.ProducerMessage
	// TODO - to bind extra metadata or log
}

func NewProducer(pc *ProducerConfig) (Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.Retry.Max = pc.Retires
	conf.Producer.RequiredAcks = sarama.WaitForAll // p.Acks
	conf.Producer.Return.Successes = true
	version, err := sarama.ParseKafkaVersion(pc.Version)
	if err != nil {
		return nil, err
	}
	conf.Version = version
	conf.ClientID = "probe" // TODO
	// conf.Metadata.Full = true
	conf.Net.SASL.Enable = true
	conf.Net.SASL.User = pc.SASLUsername
	conf.Net.SASL.Password = pc.SASLPassword
	// conf.Net.SASL.Handshake = true

	if pc.SASLMechanism == "SCRAM-SHA-512" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
	} else if pc.SASLMechanism == "SCRAM-SHA-256" {
		conf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
		conf.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
	} else {
		return nil, NoSASLMechanismFound
	}

	// TODO - tls flag - default true
	t, err := createTLSConfiguration(true, pc.SSLCAlocation, "", "")
	if err != nil {
		return nil, err
	}
	conf.Net.TLS.Enable = true
	conf.Net.TLS.Config = t

	_producer, err := sarama.NewAsyncProducer(pc.BootstrapServers, conf)
	if err != nil {
		return nil, err
	}

	p := &producer{
		producer:  _producer,
		done:      make(chan model.None),
		errors:    make(chan *model.ProducerError),
		input:     make(chan *model.ProducerMessage),
		successes: make(chan *model.ProducerMessage),
	}
	go p.loop()
	return p, nil
}

func (p *producer) Produce() chan<- *model.ProducerMessage {
	return p.input
}

func (p *producer) OnSuccess() <-chan *model.ProducerMessage {
	return p.successes
}

func (p *producer) OnError() <-chan *model.ProducerError {
	return p.errors
}

func (p *producer) Done() chan<- model.None {
	return p.done
}

func (p *producer) IsClosed() bool {
	return p.isClosed
}

func (p *producer) close() error {
	close(p.input)
	if err := p.producer.Close(); err != nil {
		return err
	}
	close(p.errors)
	close(p.successes)
	close(p.done)
	p.isClosed = true
	return nil
}

func (p *producer) toMsg(m *model.ProducerMessage) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic:     string(m.Topic),
		Key:       sarama.ByteEncoder(m.Key),
		Value:     sarama.ByteEncoder(m.Value),
		Partition: m.Partition,
	}
}

func (p *producer) fromMsg(m *sarama.ProducerMessage) *model.ProducerMessage {
	key, _ := m.Key.Encode()
	val, _ := m.Value.Encode()
	return &model.ProducerMessage{
		Topic:     model.ProducerTopic(m.Topic),
		Key:       key,
		Value:     val,
		Partition: m.Partition,
	}
}

func (p *producer) fromError(e *sarama.ProducerError) *model.ProducerError {
	return &model.ProducerError{
		Msg: p.fromMsg(e.Msg),
		Err: e.Err,
	}
}

func (p *producer) loop() {
	defer p.close()
ProducerLoop:
	for {
		select {
		case msg := <-p.input:
			p.producer.Input() <- p.toMsg(msg)
			// fmt.Println("input --", string(msg.Value))
		case err := <-p.producer.Errors():
			p.errors <- p.fromError(err)
			// fmt.Println("error --- ")
		case msg := <-p.producer.Successes():
			p.successes <- p.fromMsg(msg)
			// fmt.Println("success --", msg.Value)
		case <-p.done:
			// fmt.Println("break done --- ")
			break ProducerLoop
		}
	}
}
