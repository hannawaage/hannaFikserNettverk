package config

import (
	"time"

	. "../driver-go/elevio"
)

// Passe på at de ulike modulene ikke importerer hverandre, designe som et hieraki.

const DoorOpenTime = 3000 * time.Millisecond
const NumElevs = 3
const NumButtons = 3
const NumFloors = 4

type ElevState int

const (
	Undefined ElevState = iota - 1
	Idle
	Moving
	DoorOpen
)

type Elevator struct {
	Id     int //eller noe for å vit om master eller ikke
	Floor  int
	Dir    MotorDirection
	State  ElevState
	Orders [NumFloors][NumButtons]bool
}

type Message struct {
	Elev      Elevator
	AllOrders [NumElevs][NumFloors][NumButtons]bool
}

type EsmChns struct {
	NewOrder chan ButtonEvent
	Buttons  chan ButtonEvent
	Floors   chan int
	//Elevator 		chan Elevator
	OrderAbove chan bool
	OrderBelow chan bool
	ShouldStop chan bool
	//SignalChns chan orders.SignalChns bør ikke avhenge av orders
	//to be continued...
}
