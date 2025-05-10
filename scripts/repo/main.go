package main

import (
	"encoding/json"
	"log"
	M "probe/model"
	"probe/repository"
)

func main() {
	db, err := repository.New(1, "icmp", "")
	if err != nil {
		log.Println(err)
	}

	val := M.ICMPReq{IP: "127.0.0.1", Cid: "12345"}
	d, _ := json.Marshal(val)
	db.Create("127.0.0.1", d)

	data, err := db.Find("127.0.0.1")
	log.Println("127.0.0.1", string(data), err)

	for x := range db.Receive() {
		log.Println("curosr ---- ", string(x))
	}

}
