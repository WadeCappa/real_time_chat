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
	sockets      []*SocketData
	lastId       int
	globalLock   sync.Mutex
	socketsInUse int
}

const defaultSize = 10

func New() *EventSockets {
	sockets := EventSockets{
		sockets:      make([]*SocketData, defaultSize),
		lastId:       0,
		socketsInUse: 0,
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
	sockets.socketsInUse++
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

	// Only if we can't re-use a spot in our slice, do we append new channelb
	currentSize := len(sockets.sockets)
	sockets.doubleSizeUnsafe()

	// After doubling the size, this is going to be empty
	sockets.sockets[currentSize] = &SocketData{
		channel:  newChannel,
		uniqueId: sockets.lastId,
	}
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
			sockets.socketsInUse--
			break
		}
	}

	if sockets.socketsInUse < (len(sockets.sockets)/4) && (len(sockets.sockets)/2) >= defaultSize {
		sockets.halfSizeUnsafe()
	}
}

func (sockets *EventSockets) doubleSizeUnsafe() {
	currentSize := len(sockets.sockets)
	newSize := currentSize * 2

	for range newSize - currentSize {
		sockets.sockets = append(sockets.sockets, nil)
	}
}

func (sockets *EventSockets) halfSizeUnsafe() {
	currentSize := len(sockets.sockets)
	newSize := currentSize / 2

	newSlice := make([]*SocketData, newSize)
	count := 0
	for _, c := range sockets.sockets {
		if c != nil {
			// append to back since we usually
			// start appending from the front
			newSlice[newSize-count-1] = c
			count++
		}
	}

	fmt.Printf("moved %d sockets to the new slice of sice %d, down from %d\n", count, newSize, currentSize)
	sockets.sockets = newSlice
}
