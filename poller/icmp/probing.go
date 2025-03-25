package icmp

import (
	M "probe/model"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func ping(c M.InICMP) (out M.OutICMP, err error) {
	out = M.OutICMP{Cid: c.Cid, Params: c.Params}
	pinger, err := probing.NewPinger(c.IP)
	if err != nil {
		return
	}

	pinger.Timeout = 10 * time.Second
	if c.Params.Timeout != 0 {
		pinger.Timeout = c.Params.Timeout
	}
	pinger.Count = 4
	if c.Params.PktCount != 0 {
		pinger.Count = c.Params.PktCount
	}
	err = pinger.Run()
	if err != nil {
		return
	}

	stats := pinger.Statistics()

	// updating the output stats
	out.T = time.Now().UTC().Format("2006-01-02 15:04:05.000000")
	out.Stats = make(map[string][4]float64)
	for _, s := range c.InputStats {
		var rtt = stats.AvgRtt.Milliseconds()
		var v1 = 0
		if stats.PacketLoss != 100 {
			v1 = 100
		}
		var v2 = 100
		if rtt == -1 {
			v2 = 0
		}
		out.Stats[s.Dn] = [4]float64{float64(v1), float64(v2), stats.PacketLoss, float64(rtt)}
	}
	return
}
