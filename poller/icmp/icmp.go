package icmp

import (
	"net"
	"time"
)

type Ctx struct {
	DB      any    // DBer interface need to define
	Bucket  string // databse name mostly
	Workers int
	Timeout time.Duration
}

type C struct {
	Cid              string        `json:"cid"`
	Tag              string        `json:"tag"`
	IP               net.IPAddr    `json:"ip"`
	Timeout          time.Duration `json:"timeout"`
	Count            int           `json:"count"`
	Retries          int           `json:"retries"`
	RetryExponential bool          `json:"retry_exponential"`
}

type R struct {
	Cid        string        `json:"cid"`
	Tag        string        `json:"tag"`
	IP         net.IPAddr    `json:"ip"`
	Ok         bool          `json:"ok"`
	PacketLoss float64       `json:"packet_loss"`
	AvgRtt     time.Duration `json:"avg_rtt"`
	StdDevRtt  time.Duration `json:"std_dev_rtt"`
}

func NewICMP() *Ctx {
	return &Ctx{}
}

func (p *Ctx) Poll() error {
	return nil
}

func (p *Ctx) getPollConfig() (out []C) {
	return
}

func (p *Ctx) doPoll(_ []C, _ chan<- any) (err error) {
	return
}

func (p *Ctx) updateDB(_ R) (err error) {
	return
}

func (p *Ctx) sendToGateway(_ R) (err error) {
	// mmed to pus to call system like Kafka topic or Webhook
	return
}

func (p *Ctx) Diagnose(C) (o R) {
	return
}

func (p *Ctx) Create(ip string, config C) (err error) {
	return
}

func (p *Ctx) Delete(ip string) (err error) {
	return
}

func (p *Ctx) Get(C) R {
	return R{}
}
