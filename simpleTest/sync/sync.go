package sync

import (
	"fmt"
	"math/rand"
	"time"

	"../network/localip"
)

type Elevator struct {
	Floor int
	Dir   int
}

type Message struct {
	Elev    Elevator
	MsgId   int
	Receipt bool
	LocalIP string
	LocalID string
}

type SyncChns struct {
	SendChn     chan Message
	RecChn      chan Message
	onlineElevs chan []string
	online      chan string
}

func Sync(id string, ch SyncChns) {
	numPeers := 2 //Burde ligge i config
	elev := Elevator{3, 0}

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	go func() {
		randNr := rand.Intn(256)
		msg := Message{elev, randNr, false, localIP, id}
		for {
			ch.SendChn <- msg
			msgTimer := time.NewTimer(5 * time.Second)
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		select {
		case incomming := <-ch.RecChn:
			if id != incomming.LocalID {
				if !incomming.Receipt { //Hvis det ikke er en kvittering, skal vi svare med kvittering
					fmt.Println("Received message from ", incomming.LocalID)
					msg := Message{elev, incomming.MsgId, true, localIP, id}
					//sender ut fem kvitteringer på fem millisekunder
					for i := 0; i < 5; i++ {
						ch.SendChn <- msg
						time.Sleep(1 * time.Millisecond)
					}
				} else { //Hvis det er en kvittering, skal vi stoppe tilhørende timer
					msgTimer.Stop()
				}
			}

			if !contains(live, incomming.LocalID) { //Denne må byttes til IP når vi kjører på forskjellige
				live = append(live, incomming.LocalIP)
				ch.onlineElevs <- live
				if len(live) == numPeers {
					go func() { ch.online <- true }()
				}
			}

		case <-ch.online:
			fmt.Println("We are online")
		}
	}

	/*
		Finn ut om vi er online ved å sjekke IP på mottatte meldinger. Hvis vi har to
		forskjellige IP-er så er vi online, og da er den med lavest IP master
	*/

}

func contains(elevs []string, str string) bool {
	for _, a := range elevs {
		if a == str {
			return true
		}
	}
	return false
}
