package main

import (
	"fmt"
	"time"

	//"../network/bcast"
	"../network/peers"
)

type Message struct {
	Note string
	Id   int
}

type Receipt bool

func main() {
	id := "1"
	sendChn := make(chan bool)
	//sendChn := make(chan Message)
	//receiptChn := make(chan Receipt)
	receiptChn := make(chan peers.PeerUpdate)

	go peers.Transmitter(16571, id, sendChn)
	go peers.Receiver(16573, receiptChn)

	go func() {
		//msg := Message{"Heisann, sender mld nr ", 0}
		msg := true
		for {
			sendChn <- msg
			//msg.Id++
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case val := <-receiptChn:
			fmt.Printf("Received receipt", val)
		}
	}
}
