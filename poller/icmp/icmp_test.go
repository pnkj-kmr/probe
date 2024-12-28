package icmp_test

import (
	"encoding/json"
	M "probe/model"
	"probe/poller/icmp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	scanner := icmp.NewICMPScanner()
	assert.NotNil(t, scanner)
}

func TestScanner_Scan(t *testing.T) {
	var ping = func(c M.ICMPConf) (out M.ICMP, err error) {
		out = M.ICMP{Config: c}
		return
	}
	scanner := icmp.NewICMPScanner()
	i := M.ICMPConf{
		IP: "127.0.0.1",
	}
	icmp.Ping = ping
	out, err := scanner.Scan(i)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, i.IP, out.(M.ICMP).Config.IP)

	out, err = scanner.Scan([]byte(""))
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestPoller_Poll(t *testing.T) {
	var ping = func(c M.ICMPConf) (out M.ICMP, err error) {
		out = M.ICMP{Config: c}
		return
	}
	icmp.Ping = ping
	var db = &_mock{}

	poller := icmp.NewICMPPoller(db, "test", db, 2, 2)
	err := poller.Poll()
	assert.Nil(t, err)
}

type _mock struct{}

func (x *_mock) Find(i, j string) (d []byte, err error) {
	return
}
func (x *_mock) Cursor(i string) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		i := M.ICMPConf{
			IP: "127.0.0.1",
		}
		d, _ := json.Marshal(i)
		ch <- d
	}()
	return ch
}
func (x *_mock) Len(i string) uint64 {
	return 0
}

func (x *_mock) Receive() chan<- interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for _ = range ch {
		}
	}()
	return ch
}
