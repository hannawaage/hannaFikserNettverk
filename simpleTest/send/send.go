package main

import (
	"flag"
	"fmt"
	"time"

	"../network/bcast"
	"../network/localip"
)

type Message struct {
	Note    string
	Id      int
	Receipt bool
	LocalIP string
}

func main() {
	var id string

	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	port := 16576

	sendChn := make(chan Message)
	recChn := make(chan Message)

	go bcast.Transmitter(port, sendChn)
	go bcast.Receiver(port, recChn)

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	go func() {

		msg := Message{"hello hello, iteration is ", 0, false, localIP}
		for {
			sendChn <- msg
			msg.Id++
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-recChn:
			if !incomming.Receipt && (incomming.LocalIP == localIP) {
				fmt.Printf("Received: %#v\n", incomming.Note)
				msg := Message{"This is a receipt of message ", incomming.Id, true, localIP}
				sendChn <- msg
			} else {
				fmt.Printf("Received receipt on message %#v\n", incomming.Id, incomming.LocalIP)
			}
			/*case r := <-recReceiptCh:
			fmt.Printf("Received Receipt from", r)*/
		}
	}
}
