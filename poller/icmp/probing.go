package icmp

import (
	M "probe/model"

	probing "github.com/prometheus-community/pro-bing"
)

func ping(c M.InICMP) (out M.OutICMP, err error) {
	out = M.OutICMP{Config: c}
	pinger, err := probing.NewPinger(c.IP)
	if err != nil {
		return
	}

	pinger.Count = c.Count
	pinger.Timeout = c.Timeout
	err = pinger.Run()
	if err != nil {
		return
	}

	stats := pinger.Statistics()
	// updating the output stats
	out.Ok = stats.PacketLoss != 100
	out.PacketLoss = stats.PacketLoss
	out.AvgRtt = stats.AvgRtt
	out.StdDevRtt = stats.StdDevRtt
	return
}
