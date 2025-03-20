package exporter

import "fmt"

// Exporter
type Exporter struct {
	id   int
	name string
	// kafkaProducer any
	// webhook       any
}

func New(id int, name string) (*Exporter, error) {
	return &Exporter{id: id, name: name}, nil
}

func (r *Exporter) Produce(data []byte) error {
	// in case of kafka
	// - need to group the extra metadata into message
	// - multi message should be group together to push to kafka
	// - it will help to partition based topic push
	// - mode information needed like org, agent_id or etc
	fmt.Println("kafka push------", r.name, data)
	return nil
}
