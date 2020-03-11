package main

import (
	"time"

	"../network/bcast"
)

type Message struct {
	Note string
	Id   int
}

func main() {
	sendChn := make(chan Message)

	go bcast.Transmitter(16571, sendChn)
	go func() {
		msg := Message{"Heisann, sender mld nr ", 0}
		for {
			sendChn <- msg
			msg.Id++
			time.Sleep(1 * time.Second)
		}
	}()
	for {
		/*
			select {
			case outgoing := <-sendChn:
				fmt.Printf("Prepared for sending: %#v\n", outgoing)
			}*/
	}
}
