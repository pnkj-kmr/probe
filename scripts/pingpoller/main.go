package main

import (
	M "probe/model"
	"probe/poller"
	"probe/poller/icmp"
	"probe/repository/mockdata"
	"time"
)

func main() {
	var ping = func(c M.ICMPReq) (out M.ICMPRes, err error) {
		out = M.ICMPRes{Cid: c.Cid, Params: c.Params}
		return
	}
	icmp.Ping = ping
	var cp = mockdata.NewPingMock()

	poller := icmp.NewPoller(cp[0], nil, poller.WithWorkers(2), poller.WithTimeout(2*time.Second))
	poller.Poll()

	// fmt.Println("last error --", err)
}
