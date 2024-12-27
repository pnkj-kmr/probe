package icmp

import (
	probing "github.com/prometheus-community/pro-bing"
)

func ping(c C) (out R, err error) {
	pinger, err := probing.NewPinger(c.IP.String())
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
	out = R{
		Cid:        c.Cid,
		Tag:        c.Tag,
		IP:         c.IP,
		Ok:         stats.PacketLoss != 100,
		PacketLoss: stats.PacketLoss,
		AvgRtt:     stats.AvgRtt,
		StdDevRtt:  stats.StdDevRtt,
	}
	return
}
