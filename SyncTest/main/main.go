package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	//. "../config"
	//. "../driver-go/elevio"
	"../network/bcast"
	"../network/localip"
	"../network/peers"
)

const numFloors = 4

// HelloMsg We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
type HelloMsg struct {
	Message string
	Iter    int
}

func main() {
	//Init("localhost:22222", NumFloors)
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	// ... or alternatively, we can use the local IP address.
	// (But since we can run multiple programs on the same PC, we also append the
	//  process ID)

	if id == "" {
		localIP, err := localip.LocalIP()
		if err != nil {
			fmt.Println(err)
			localIP = "DISCONNECTED"
		}
		id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())
	}

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15648, id, peerTxEnable)
	go peers.Receiver(15648, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan HelloMsg)
	helloRx := make(chan HelloMsg)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16570, helloTx)
	go bcast.Receiver(16570, helloRx)

	// The example message. We just send one of these every second.
	go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-helloRx:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}

/*import (
	"flag"
	"fmt"
	"time"

	. "../config"
	. "../driver-go/elevio"
	"../network/bcast"
)

type iAmAlive struct {
	Id    string
	Floor int
}

type SyncChns struct {
	Status chan iAmAlive
	Update chan iAmAlive
	//to be continued...
}

func main() {

	syncChns := SyncChns{
		Status: make(chan iAmAlive),
		Update: make(chan iAmAlive),
	}

	var id string
	flag.StringVar(&id, "id", "", "id of this peer")
	flag.Parse()

	Init("localhost:22222", NumFloors)
	elevator := Elevator{
		Id:     1, //eller noe for å vit om master eller ikke
		Floor:  0,
		Dir:    MD_Up,
		State:  Idle,
		Orders: [NumFloors][NumButtons]bool{},
	}

	go bcast.Transmitter(16569, syncChns.Update)
	go bcast.Receiver(16569, syncChns.Status)

	go func() {
		msg := iAmAlive{id, elevator.Floor}
		for {
			msg.Floor++
			syncChns.Update <- msg
			time.Sleep(1 * time.Second)
		}
	}()

	fmt.Println("Started")

	for {
		select {
		case u := <-syncChns.Update:
			fmt.Printf("Update transmitted by :\n", u.Id)
		case f := <-syncChns.Status:
			fmt.Printf("Status recieved by :\n", f.Id)
		}
	}

}
*/
