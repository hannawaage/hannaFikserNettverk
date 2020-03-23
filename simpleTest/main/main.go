package main

import (
	"flag"

	"../network/bcast"
	"../sync"
)

func main() {
	syncChns := sync.SyncChns{
		SendChn: make(chan sync.Message),
		RecChn:  make(chan sync.Message),
		Online:  make(chan bool),
	}

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	port := 16576

	go bcast.Transmitter(port, syncChns.SendChn)
	go bcast.Receiver(port, syncChns.RecChn)
	go sync.Sync(id, syncChns)
	go sync.OrdersDist(syncChns)
	for {

	}
}
