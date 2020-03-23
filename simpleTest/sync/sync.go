package sync

import (
	"fmt"
	"math/rand"
	"strconv"
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
	SendChn chan Message
	RecChn  chan Message
	Online  chan bool
}

func Sync(id string, ch SyncChns) {
	const numPeers = 2 //Burde ligge i config
	elev := Elevator{3, 0}

	var (
		iAmMaster bool = true
		//online    bool //boolean initiates to false
		onlineIPs       []string
		receivedReceipt []string
		currentMsgID    int
	)

	localIP, err := localip.LocalIP()
	if err != nil {
		fmt.Println(err)
		localIP = "DISCONNECTED"
	}

	msgTimer := time.NewTimer(5 * time.Second)
	msgTimer.Stop()

	go func() {
		currentMsgID = rand.Intn(256)
		msg := Message{elev, currentMsgID, false, localIP, id}
		for {
			ch.SendChn <- msg
			msgTimer.Reset(1 * time.Second)
			time.Sleep(1 * time.Second)
			/*
				go func() {
					msgRec := <-ch.timerConf
					if msgRec == randNr {
						msgTimer.Stop()
					}
				}()
			*/
		}
	}()

	for {
		select {
		case incomming := <-ch.RecChn:
			recID := incomming.LocalID
			if id != recID { //Hvis det ikke er fra oss selv, BYTTES TIL IP VED KJØRING PÅ FORSKJELLIGE MASKINER
				if !contains(onlineIPs, recID) {
					// Dersom heisen enda ikke er registrert, sjekker vi om vi nå er online og sjekker om vi er master
					onlineIPs = append(onlineIPs, recID)
					if len(onlineIPs) == numPeers {
						ch.Online <- true
						idDig, _ := strconv.Atoi(id)
						for i := 0; i < numPeers; i++ {
							theID, _ := strconv.Atoi(onlineIPs[i])
							if idDig > theID {
								iAmMaster = false
								break
							}
						}
						/*
							Dette er ved diff på IP:
							localDig, _ := strconv.Atoi(localIP[len(localIP)-3:])
							for i := 0; i <= numPeers; i++ {
								theIP := onlineIPs[i]
								lastDig, _ := strconv.Atoi(theIP[len(theIP)-3:])
								if localDig < lastDig {
									iAmMaster = false
									break
								}
							}
						*/
						if iAmMaster {
							fmt.Println("I am master")
						} else {
							fmt.Println("I am backup")
						}
					}
				}
				if !incomming.Receipt {
					// Hvis det ikke er en kvittering, skal vi svare med kvittering
					msg := Message{elev, incomming.MsgId, true, localIP, id}
					//sender ut fem kvitteringer på femti millisekunder
					for i := 0; i < 5; i++ {
						ch.SendChn <- msg
						time.Sleep(10 * time.Millisecond)
					}
				} else { // Hvis det er en kvittering
					if incomming.MsgId == currentMsgID {
						if !contains(receivedReceipt, recID) {
							receivedReceipt = append(receivedReceipt, recID)
							if len(receivedReceipt) == numPeers {
								msgTimer.Stop()
								receivedReceipt = receivedReceipt[:0]
							}
						}
					}
				}
			}
		case <-msgTimer.C:
			ch.Online <- false
			fmt.Println("Timer timeout")
		}
	}
}

func OrdersDist(ch SyncChns) {
	var (
		online bool //initiates to false
	)
	go func() {
		for {
			select {
			case b := <-ch.Online:
				if b {
					online = true
					fmt.Println("Yaho, we are online!")
				} else {
					online = false
					fmt.Println("Boo, we are offline.")
				}
			}
		}
	}()
	for {
		if online {
			fmt.Println("Online")
			time.Sleep(5 * time.Second)
		}
	}

}

func contains(elevs []string, str string) bool {
	for _, a := range elevs {
		if a == str {
			return true
		}
	}
	return false
}
