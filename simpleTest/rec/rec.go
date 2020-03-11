package main

import (
	"fmt"

	//"../network/bcast"
	"../network/peers"
)

type Message struct {
	Note string
	Id   int
}

type Receipt bool

func main() {
	id = "2"
	recChn := make(chan bool)
	receiptChn := make(chan peers.PeerUpdate)
	//recChn := make(chan Message)
	//receiptChn := make(chan Receipt)

	go peers.Receiver(16571, recChn)
	go peers.Transmitter(16573, id, receiptChn)

	fmt.Printf("heihei")
	myPeer := peers.PeerUpdate{"", "2", ""}
	for {
		select {
		case incomming := <-recChn:
			fmt.Printf("Received: %#v\n", incomming)
			go func() {
				receiptChn <- true
				fmt.Printf("Receipt sent")
			}()
		}
	}

}
