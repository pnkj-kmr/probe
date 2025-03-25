package icmp_test

import (
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"
	"probe/repository/mockdata"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	var cp = mockdata.NewPingMock()
	scanner := icmp.NewScanner(cp, cp, poller.WithWorkers(2), poller.WithTimeout(2*time.Second))
	assert.NotNil(t, scanner)
}

func TestScanner_Scan(t *testing.T) {
	var ping = func(c M.InICMP) (out M.OutICMP, err error) {
		out = M.OutICMP{Params: c.Params, Cid: c.Cid}
		return
	}
	var cp = mockdata.NewPingMock()
	scanner := icmp.NewScanner(cp, cp, poller.WithWorkers(2), poller.WithTimeout(2*time.Second))
	i := M.InICMP{
		IP: "127.0.0.1",
	}
	icmp.Ping = ping
	out, err := scanner.Scan(i)
	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, i.Cid, out.(M.OutICMP).Cid)

	out, err = scanner.Scan([]byte(""))
	assert.Nil(t, out)
	assert.NotNil(t, err)
}

func TestPoller_Poll(t *testing.T) {
	var ping = func(c M.InICMP) (out M.OutICMP, err error) {
		out = M.OutICMP{Params: c.Params, Cid: c.Cid}
		return
	}
	icmp.Ping = ping
	var cp = mockdata.NewPingMock()

	p := icmp.NewPoller(cp, cp, poller.WithWorkers(2), poller.WithTimeout(2*time.Second))
	err := p.Poll()
	assert.Nil(t, err)

	var ping2 = func(c M.InICMP) (out M.OutICMP, err error) {
		return out, icmp.ErrICMP
	}
	icmp.Ping = ping2
	var cp2 = mockdata.NewPingMock()
	p = icmp.NewPoller(cp2, cp2, poller.WithWorkers(2), poller.WithTimeout(2*time.Second))
	err = p.Poll()
	assert.Nil(t, err)
	assert.NotEmpty(t, cp2.Count)
}
