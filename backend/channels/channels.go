package channels

import (
	"fmt"
	"sync"
)

type SocketData struct {
	uniqueId int
	channel  chan []byte
}
type EventSockets struct {
	sockets    []*SocketData
	lastId     int
	globalLock sync.Mutex
}

func New() *EventSockets {
	defaultSize := 10
	sockets := EventSockets{
		sockets: make([]*SocketData, defaultSize),
		lastId:  0,
	}

	return &sockets
}

func (sockets *EventSockets) FanInMessage(event []byte) {
	eventCopy := event
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := sockets.sockets[i]
		if c != nil {
			fmt.Printf("talking to socket %d\n", c.uniqueId)
			c.channel <- eventCopy
		}
	}
}

func (sockets *EventSockets) AddChannel(newChannel chan []byte) int {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	sockets.lastId++
	defer fmt.Printf("created channel of id %d\n", sockets.lastId)
	for i := range sockets.sockets {
		c := sockets.sockets[i]
		if c == nil {
			newChannel := SocketData{
				channel:  newChannel,
				uniqueId: sockets.lastId,
			}
			sockets.sockets[i] = &newChannel
			return sockets.lastId
		}
	}

	// Only if we can't re-use a spot in our slice, do we append a new channel
	sockets.sockets = append(sockets.sockets, &SocketData{
		channel:  newChannel,
		uniqueId: sockets.lastId,
	})
	return sockets.lastId
}

func (sockets *EventSockets) RemoveChannel(id int) {
	sockets.globalLock.Lock()
	defer sockets.globalLock.Unlock()
	for i := range sockets.sockets {
		c := sockets.sockets[i]
		if c != nil && c.uniqueId == id {
			close(c.channel)
			sockets.sockets[i] = nil
			fmt.Printf("removed socket %d\n", id)
			return
		}
	}
	fmt.Printf("Could not find socket %d\n", id)
}
