package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"probe/repository"
	"strconv"
	"time"
)

func main() {
	st := time.Now()
	c := 1000
	db, err := repository.New(fmt.Sprintf("icmp_%d", c), 0)
	if err != nil {
		log.Println(err)
	}

	// //inserting
	sample_data := []byte(`{"ci_id": "125654075603791384629", "ptype": "ping","ip": "99.99.99.99","ip_version": 4,"login_profileid": "123153408056393469952","business_profileid": 0, "poll_period": 60,"input_stats": [{ "dn": "input", "oid": "1" }],"params": { "pollPeriod": 60, "timeout": "10", "packsiz": "32", "interval": "60", "poll_addr": "99.99.99.99", "poll_period": 60, "configured_speed": 2000000.0, "secondary_poll_addr": "99.99.99.99", "wan_ip_ping_check": 1, "wan_ip_addr": "99.99.99.99", "mib_profile": "ping.cfg", "ci_name": "ECLT-BRR128-RTR3614_99.99.99.99", "ci_category": "interface", "parent_ci_id": "125654075603791384625", "poll_type": "ping"}}`)
	fmt.Println("inserting into db with records-", c)
	var k string
	for i := 0; i < c; i++ {
		k = strconv.Itoa(rand.Intn(math.MaxInt))
		db.Create(k, sample_data)
	}
	fmt.Println("Insert : time", time.Since(st))
	st = time.Now()

	_c := 0
	for range <-db.Receive() {
		_c = _c + 1
		// fmt.Println(i, _c)

	}
	fmt.Println("totale records --", _c)
	fmt.Println("Fetch : time", time.Since(st))

}
