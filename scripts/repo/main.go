package main

import (
	"log"
	"probe/repository"
)

func main() {
	db, err := repository.New("icmp", 1)
	if err != nil {
		log.Println(err)
	}

	// val := []byte("v1")
	// db.Create("k1", val)

	data, err := db.Find("k1")
	log.Println("k1", string(data), err)

	for x := range db.Consume() {
		log.Println("curosr ---- ", string(x))
	}

}
