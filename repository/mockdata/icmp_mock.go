package mockdata

import (
	"encoding/json"
	"fmt"
	M "probe/model"
	"sync"
)

type PingMock struct {
	mu    sync.Mutex
	Count int
	ch    chan []byte
}

func NewPingMock() *PingMock {
	p := &PingMock{ch: make(chan []byte)}
	go p.spin()
	return p
}

func (x *PingMock) Receive() <-chan []byte {
	return x.ch
}

func (x *PingMock) spin() {
	for k := 1; k < 9; k++ {
		if k == 4 {
			i := M.ICMPRes{
				Cid:   "44444",
				Stats: map[string][4]float64{"input": {1, 2, 3, 4}},
			}
			d, _ := json.Marshal(i)
			x.ch <- d
		}
		if k == 3 {
			i := M.ICMPRes{
				Cid: "12345",
			}
			d, _ := json.Marshal(i)
			x.ch <- d
		}
		if k == 2 {
			i := M.ICMPReq{
				IP: "127.0.0.2",
			}
			d, _ := json.Marshal(i)
			x.ch <- d
		}
		if k == 1 {
			i := M.ICMPReq{
				IP: "127.0.0.1",
			}
			d, _ := json.Marshal(i)
			x.ch <- d
		}
		if k == 5 {
			close(x.ch)
		}
	}

}

func (x *PingMock) Send(d any) (err error) {
	x.mu.Lock()
	x.Count = x.Count + 1
	x.mu.Unlock()
	fmt.Println("-------- data produced ---", x.Count)
	return
}
