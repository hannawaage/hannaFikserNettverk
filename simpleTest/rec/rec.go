package main

import (
	"fmt"

	"../network/bcast"
)

type Message struct {
	Note string
	Id   int
}

func main() {
	recChn := make(chan Message)

	go bcast.Receiver(16571, recChn)

	for {
		select {
		case incomming := <-recChn:
			fmt.Printf("Received: %#v\n", incomming)
		}
	}

}
